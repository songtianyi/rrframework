package rrserver

import (
	"fmt"
	"io"
	"net"
	"strconv"
)

var (
	CustomHandleConn = func(c *TCPConnection, packet []byte) { fmt.Println("forget rrserver.CustomHandleConn = YourHanleConn in init func?") }
)

type TCPServer struct {
	ls   net.Listener
	port int
}

func CreateTCPServer(inf string, port int) (error, *TCPServer) {
	err, ipaddr := getIpAddrByInterface(inf)
	if err != nil {
		return err, nil
	}
	listener, err := net.Listen("tcp", net.JoinHostPort(ipaddr, strconv.Itoa(port)))
	if err != nil {
		return err, nil
	}
	s := &TCPServer{
		ls:   listener,
		port: port,
	}
	return nil, s
}

func (s *TCPServer) Start() {
	fmt.Printf("Server listening in [%s]\n", s.ls.Addr())
	for {
		conn, err := s.ls.Accept()
		if err != nil {
			fmt.Println("Server Accept() return error, %s", err)
			break
		}
		fmt.Printf("new msg [%s]-->[%s]\n", conn.RemoteAddr(), conn.LocalAddr())
		go s.handleConn(NewTCPConnection(conn)) 
	}
	return
}

func (s *TCPServer) handleConn(c *TCPConnection) {
	for {
		err, packet := c.Read()
		if err != nil {
			// end goroutine
			if err != io.EOF {
				fmt.Println(err)
				return
			}else{
				fmt.Println("EOF")
				return
			}
		}
		go CustomHandleConn(c, packet)
	}
}
