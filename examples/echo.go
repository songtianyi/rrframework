package main

import (
	"fmt"
	"time"
	"rrframework/handler"
	"rrframework/server"
//	"./proto/example"
	"./proto/rrfp"
)

func echo(c interface{}, msg interface{}) {
	conn := c.(*rrserver.TCPConnection)
	m := msg.(*rrfp.Message)
	fmt.Println(conn, 1)
	fmt.Println(m.String())
}

func main() {
	err, hr := rrhandler.CreateHandleRegister()
	if err != nil {
		fmt.Println(err)
	}
	hr.Add("example.EchoRequest", rrhandler.HandlerFuncWrapper(echo), 0*time.Second)

	//
	err, s := rrserver.CreateTCPServer("eth0", 8003, hr)
	if err != nil {
		fmt.Println(err)
		return
	}
	s.Start()
}
