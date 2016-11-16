package main

import (
	"./proto/rrfp"
	"github.com/golang/protobuf/proto"
	"rrframework/handler"
	"rrframework/logs"
	"rrframework/server"
	"rrframework/utils"
	"time"
)

var (
	hr *rrhandler.HandlerRegister
)

func echo(c interface{}, msg interface{}) {
	conn := c.(*rrserver.TCPConnection)
	m := msg.(*rrfp.Message)
	logs.Debug("Request msg:", m.GetBy().GetExampleEchoRequest().Msg)

	req := new(rrfp.Message)
	req.Hd = &rrfp.Head{
		rrutils.NewV4().String(),
		"rrfp.ExampleEchoResponse",
	}
	req.By = &rrfp.Body{
		MsgType: &rrfp.Body_ExampleEchoResponse{
			ExampleEchoResponse: &rrfp.ExampleEchoResponse{Msg: "Lucky!"},
		},
	}
	b, _ := proto.Marshal(req)

	if err := conn.Write(b); err != nil {
		logs.Error(err)
		return
	}
	return
}

func init() {
	_, hr = rrhandler.CreateHandlerRegister()
	hr.Add("rrfp.ExampleEchoRequest", rrhandler.Handler(echo), 0*time.Second)

	rrserver.CustomHandleConn = HandleConn
}

func HandleConn(c *rrserver.TCPConnection, packet []byte) {
	logs.Debug("new msg [%s]-->[%s]", c.RemoteAddr(), c.LocalAddr())
	msg := new(rrfp.Message)
	err := proto.Unmarshal(packet, msg)
	if err != nil {
		logs.Debug("Unmarshal packet err, %s", err)
		return
	}
	err, hw := hr.Get(msg.GetHd().UniqueId)
	if err != nil {
		logs.Debug("Can't find handle for message type [%s], %s", msg.GetHd().UniqueId, err)
		return
	}
	go hw.Run(c, msg)
}

func main() {

	err, s := rrserver.CreateTCPServer("0.0.0.0", 8003)
	if err != nil {
		logs.Debug(err)
		return
	}
	s.Start()
}
