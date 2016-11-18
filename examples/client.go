package main

import (
	"github.com/golang/protobuf/proto"
	"net"
	"rrframework/examples/proto/rrfp"
	"rrframework/logs"
	"rrframework/server"
	"rrframework/utils"
)

func main() {
	conn, err := net.Dial("tcp", "0.0.0.0:8003")
	if err != nil {
		logs.Error(err)
		return
	}
	c := rrserver.NewTCPConnection(conn)

	msg := new(rrfp.Message)
	msg.Hd = &rrfp.Head{
		rrutils.NewV4().String(),
		"rrfp.ExampleEchoRequest",
	}
	msg.By = &rrfp.Body{
		MsgType: &rrfp.Body_ExampleEchoRequest{
			ExampleEchoRequest: &rrfp.ExampleEchoRequest{Msg: "hello world!"},
		},
	}
	logs.Debug("before marshal:", msg)
	b, err := proto.Marshal(msg)
	if err != nil {
		logs.Error(err)
		return
	}

	if err := c.Write(b); err != nil {
		logs.Error(err)
		return
	}

	err, packet := c.Read()
	if err != nil {
		logs.Error(err)
		return
	}
	m := new(rrfp.Message)
	proto.Unmarshal(packet, m)
	logs.Info(m.String())
	logs.Debug("Response msg", m.GetBy().GetExampleEchoResponse().Msg)

}
