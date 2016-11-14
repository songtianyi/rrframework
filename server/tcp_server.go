package rrserver

import (
	"fmt"
	"io"
	"net"
	"rrframework/handler"
	"rrframework/examples/proto/rrfp"
	"rrframework/serializer/protobuf/proto"
	"strconv"
)

type TCPServer struct {
	ls   net.Listener
	port int
	hr   *rrhandler.HandleRegister
}

func CreateTCPServer(inf string, port int, hr *rrhandler.HandleRegister) (error, *TCPServer) {
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
		hr:   hr,
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
		defer c.conn.Close()
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
				return
			}else{
				return
			}
		}
		msg := new(rrfp.Message)
		err = proto.Unmarshal(packet, msg)
		fmt.Println("got message")
		if err != nil {
			fmt.Sprintf("Unmarshal packet err, %s", err)
			continue
		}
		err, handle := s.hr.Get(msg.GetHd().UniqueId)
		fmt.Println(handle)
		if err != nil {
			fmt.Println("handle for 1 not registred")
			continue
		}
		go handle.(rrhandler.Handler).Run(c, msg)
	}
}
