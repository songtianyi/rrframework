package main

import (
	"github.com/golang/protobuf/proto"
	"rrframework/connector/redis"
	"rrframework/connector/zookeeper"
	"rrframework/examples/proto/rrfp"
	"rrframework/handler"
	"rrframework/logs"
	"rrframework/server"
	"rrframework/utils"
	"time"
)

var (
	hr *rrhandler.HandlerRegister
)

func joke(msg interface{}) {
	logs.Debug("joke")
	m := msg.(*rrfp.Message)
	b, err := proto.Marshal(m)
	if err != nil {
		logs.Error(err)
	}
	// send the same msg to itself
	err, newRes := rrserver.SendTCPRequest("127.0.0.1:8003", b)
	logs.Debug(err, newRes)
}

func echo(c interface{}, msg interface{}) {
	conn := c.(*rrserver.TCPConnection)
	m := msg.(*rrfp.Message)
	logs.Debug("Request msg:", m.GetBy().GetExampleEchoRequest().Msg)

	res := new(rrfp.Message)
	res.Hd = &rrfp.Head{
		rrutils.NewV4().String(),
		"rrfp.ExampleEchoResponse",
	}
	res.By = &rrfp.Body{
		MsgType: &rrfp.Body_ExampleEchoResponse{
			ExampleEchoResponse: &rrfp.ExampleEchoResponse{Msg: "Lucky!"},
		},
	}

	// connect redis, set msg to db
	err, rc := rrredis.GetRedisClient("127.0.0.1:6379")
	if err != nil {
		logs.Error(err)
	} else {
		result, err := rc.Get("songtianyi")
		if err != nil {
			logs.Debug(err)
			res.GetBy().GetExampleEchoResponse().Msg = err.Error()
		} else {
			res.GetBy().GetExampleEchoResponse().Msg = string(result)
		}
	}

	b, _ := proto.Marshal(res)
	if err := conn.Write(b); err != nil {
		logs.Error(err)
		return
	}

	//joke(msg)
	err, _ = rrzk.GetZkClient("10.19.150.38:2181,10.19.168.143:2181,10.19.3.141:2181")
	if err != nil {
		logs.Error(err)
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
