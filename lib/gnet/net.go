package gnet

import (
	"net"

	"github.com/u35s/gmod/lib/utils"
)

func Listen(addr string) (*net.TCPListener, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		utils.Err("resolve tcp addr  %v err,%v", addr, err)
		return nil, err
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		utils.Err("listen tcp %v err,%v", addr, err)
		return nil, err
	}
	utils.Inf("listen %v success", addr)
	return listener, nil
}

func Accept(listener *net.TCPListener, f func(conn net.Conn)) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			utils.Err("accept listener %v err,%v", listener.Addr(), err)
			break
		}
		go f(conn)
	}
}
