package gnet

import (
	"errors"
	"net"
	"time"

	"github.com/u35s/gmod/lib/utils"
)

func ConnectTo(addr string) (*net.TCPConn, error) {
	c, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		return nil, err
	}
	conn, ok := c.(*net.TCPConn)
	if !ok {
		return nil, errors.New("conn to tcpconn err")
	}
	conn.SetKeepAlivePeriod(time.Second)
	conn.SetKeepAlive(true)
	return conn, nil
}

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
		conn, err := listener.AcceptTCP()
		if err != nil {
			utils.Err("accept listener %v err,%v", listener.Addr(), err)
			break
		}
		conn.SetKeepAlivePeriod(time.Second)
		conn.SetKeepAlive(true)
		go f(conn)
	}
}
