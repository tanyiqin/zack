package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/tanyiqin/zack/pb"
	"github.com/tanyiqin/zack/znet"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("conn err,", err)
		return
	}
	msg := &pb.CsAccountLogin{RoleID: 123456, PassWord: "ass11www"}
	data, err := proto.Marshal(msg)
	if err != nil {
		fmt.Println("proto marshal err", err)
		return
	}

	msg1, err := znet.Pack(znet.NewMessage(1, data))
	if err != nil {
		fmt.Println("pack err", err)
		return
	}

	_, err = conn.Write(msg1)
	if err != nil {
		fmt.Println("write err", err)
		return
	}

	msg2, _ := znet.ReadFromConn(conn)

	result1 := &pb.ScAccountLoginResult{}
	proto.Unmarshal(msg2.GetMsgData(), result1)

	fmt.Println("return result=", result1.Result)

}
