package gnet

import (
	"net"
	"time"

	"github.com/u35s/gmod/lib/utils"
	"github.com/u35s/rudp"
)

func Dial(network, addr string) (net.Conn, error) {
	if network == "udp" {
		laddr := net.UDPAddr{IP: net.IPv4zero, Port: 0}
		conn, err := net.DialUDP("udp", &laddr, getUDPAddr(addr))
		if err != nil {
			return nil, err
		}
		return rudp.NewConn(conn, rudp.New()), nil
	}

	return net.DialTimeout(network, addr, 5*time.Second)
}

func getUDPAddr(addr string) *net.UDPAddr {
	ip := "0.0.0.0"
	port := ""
	slc := utils.Split(addr, ":")
	if len(slc) == 1 {
		port = slc[0]
	} else if len(slc) == 2 {
		ip = slc[0]
		port = slc[1]
	}
	return &net.UDPAddr{IP: net.ParseIP(ip), Port: utils.Atoi(port)}
}

func ConnectTo(network, addr string) (net.Conn, error) {
	if network == "udp" {
		laddr := net.UDPAddr{IP: net.IPv4zero, Port: 0}
		conn, err := net.DialUDP("udp", &laddr, getUDPAddr(addr))
		if err != nil {
			return nil, err
		}
		return rudp.NewConn(conn, rudp.New()), nil
	}
	c, err := net.DialTimeout(network, addr, 1*time.Second)
	if err != nil {
		return nil, err
	}
	if tcp, ok := c.(*net.TCPConn); ok {
		tcp.SetKeepAlivePeriod(time.Second)
		tcp.SetKeepAlive(true)
	}
	return c, nil
}

func Listen(network, addr string) (net.Listener, error) {
	check := func(err error) bool {
		if err != nil {
			utils.Err("listen %v %v err,%v", network, addr, err)
			return true
		}
		utils.Inf("listen %v %v success", network, addr)
		return false
	}
	if network == "udp" {
		conn, err := net.ListenUDP("udp", getUDPAddr(addr))
		if check(err) {
			return nil, err
		}
		return rudp.NewListener(conn), nil
	}
	listener, err := net.Listen(network, addr)
	if check(err) {
		return nil, err
	}
	return listener, nil
}

func Accept(listener net.Listener, f func(conn net.Conn)) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			utils.Err("accept listener %v err,%v", listener.Addr(), err)
			break
		}
		if tcp, ok := conn.(*net.TCPConn); ok {
			tcp.SetKeepAlivePeriod(time.Second)
			tcp.SetKeepAlive(true)
		}
		go f(conn)
	}
}
