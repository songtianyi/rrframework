package main

import (
	"fmt"
	"net"
	"rrframework/server"
	"rrframework/utils"
	"rrframework/examples/proto/example"
	"./proto/rrfp"
	"rrframework/serializer/protobuf/proto"
)

func main() {
	conn, err := net.Dial("tcp", "10.19.147.75:8003")
	if err != nil {
		fmt.Println(err)
		return
	}
	c := rrserver.NewTCPConnection(conn)
	fmt.Println(c)
	msg := new(rrfp.Message)
	msg.Hd = &rrfp.Head{
		rrutils.NewV4().String(),	
		"example.EchoRequest",
	}
	msg.By = &rrfp.Body{
		&rrfp.Body_ExampleEchoRequest{
			&example.EchoRequest{
				"fuck you man",
			},
		},
	}
	b, _ := proto.Marshal(msg)
	c.Write(b)

}
