package main

import (
	"./proto/rrfp"
	"fmt"
	"rrframework/serializer/protobuf/proto"
	"rrframework/utils"
)

func main() {
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
	fmt.Println("before marshal:", msg)
	b, err := proto.Marshal(msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	// unmarshal
	m := new(rrfp.Message)
	if err := proto.Unmarshal(b, m); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("after unmarshal:", m)
	fmt.Println("m.GetBy()", m.GetBy())
	fmt.Println("m.GetBy().GetMsgType()", m.GetBy().GetMsgType())
	fmt.Println("m.GetBy().GetExampleEchoRequest()", m.GetBy().GetExampleEchoRequest())
}
