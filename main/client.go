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
	// 注册
	msg := &pb.CsPlayerCreate{Name: "wewe"}
	data, _ := proto.Marshal(msg)
	msg1, err := znet.Pack(znet.NewMessage(2, data))
	_, err = conn.Write(msg1)
	msg2, _ := znet.ReadFromConn(conn)
	result1 := &pb.ScPlayerCreateResult{}
	proto.Unmarshal(msg2.GetMsgData(), result1)
	fmt.Println(result1.Result)

	// 登录
	//msg := &pb.CsAccountLogin{RoleID: 123456}
	//data, err := proto.Marshal(msg)
	//if err != nil {
	//	fmt.Println("proto marshal err", err)
	//	return
	//}
	//
	//msg1, err := znet.Pack(znet.NewMessage(1, data))
	//if err != nil {
	//	fmt.Println("pack err", err)
	//	return
	//}
	//
	//_, err = conn.Write(msg1)
	//if err != nil {
	//	fmt.Println("write err", err)
	//	return
	//}
	//
	//msg2, _ := znet.ReadFromConn(conn)
	//
	//result1 := &pb.ScAccountLoginResult{}
	//proto.Unmarshal(msg2.GetMsgData(), result1)
	//
	//fmt.Println("return result=", result1.Result)

}
