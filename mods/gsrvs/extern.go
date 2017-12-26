package gsrvs

import (
	"log"
)

var srvs *ConnectedServerSlcSafeMap

func init() {
	srvs = new(ConnectedServerSlcSafeMap)
	srvs.init()
}

func SendCmdToServerWithID(id uint, msg interface{}) {
	srvs.lock.RLock()
	for _, slc := range srvs.connectedServerSlcMap {
		for _, v := range slc {
			if id == 0 || v.ID == id {
				v.Agent.SendCmd(msg)
			}
		}
	}
	srvs.lock.RUnlock()
}

func SendCmdToServer(tp, name string, msg interface{}) {
	srvs.lock.RLock()
	if slc, ok := srvs.connectedServerSlcMap[tp]; ok {
		for _, v := range slc {
			if name == "" || v.Name == name {
				v.Agent.SendCmd(msg)
			}
		}
	}
	srvs.lock.RUnlock()
}

func AddToListenAddr(addr, net string) {
	srvs.toListen = append(srvs.toListen, &ToListenAddr{Addr: addr, Net: net})
}

func EachToListenAddr(f func(*ToListenAddr)) {
	for _, v := range srvs.toListen {
		f(v)
	}
}

func AddToConnectServerWithID(id uint, tp, name, addr, net string) {
	srvs.toConnect = append(srvs.toConnect, &ToConnectServer{ServerBase: ServerBase{Type: tp, Name: name, ID: id},
		Addr: addr, Net: net})
}

func AddToConnectServer(tp, name, addr, net string) {
	AddToConnectServerWithID(0, tp, name, addr, net)
}

func EachToConnectServer(f func(*ToConnectServer)) {
	for _, v := range srvs.toConnect {
		f(v)
	}
}

func Has(tp, name string) bool {
	srvs.lock.RLock()
	defer srvs.lock.RUnlock()
	slc, ok := srvs.connectedServerSlcMap[tp]
	if ok {
		for i := range slc {
			if slc[i].Name == name {
				return true
			}
		}
	}
	return false
}

func Add(srv *ConnectedServer) {
	srvs.lock.Lock()
	slc, ok := srvs.connectedServerSlcMap[srv.Type]
	if !ok {
		slc = make(ConnectedServerSlc, 0, 1)
	}
	slc = append(slc, srv)
	srvs.connectedServerSlcMap[srv.Type] = slc
	srvs.size++
	srvs.lock.Unlock()
	log.Printf("%v,%v add success,local addr %v,remote addr %v",
		srv.Type, srv.Name, srv.Agent.Conn.LocalAddr(), srv.Agent.Conn.RemoteAddr())
}

func Remove(srv *ConnectedServer) {
	srvs.lock.Lock()
	if slc, ok := srvs.connectedServerSlcMap[srv.Type]; ok {
		for i := range slc {
			if slc[i] == srv {
				slc = append(slc[:i], slc[i+1:]...)
				srvs.connectedServerSlcMap[srv.Type] = slc
				srvs.size--
				log.Printf("%v,%v remove success,local addr %v,remote addr %v,index %v",
					srv.Type, srv.Name, srv.Agent.Conn.LocalAddr(), srv.Agent.Conn.RemoteAddr(), i)
				EachToConnectServer(func(to *ToConnectServer) {
					if to.Type == srv.Type && to.Name == srv.Name && to.ID == srv.ID {
						to.Ok = false
					}
				})
				break
			}
		}

	}
	srvs.lock.Unlock()
}

func TypeSize() int   { return len(srvs.connectedServerSlcMap) }
func ServerSize() int { return srvs.size }
