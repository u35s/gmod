package gnet

import (
	"log"
	"net"
)

func Listen(addr string) (*net.TCPListener, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		log.Printf("resolve tcp addr  %v err,%v", addr, err)
		return nil, err
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Printf("listen tcp %v err,%v", addr, err)
		return nil, err
	}
	log.Printf("listen %v success", addr)
	return listener, nil
}

func Accept(listener *net.TCPListener, f func(conn net.Conn)) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept listener %v err,%v", listener.Addr(), err)
			break
		}
		go f(conn)
	}
}
