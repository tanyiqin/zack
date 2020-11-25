package router

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/tanyiqin/zack/pb"
	"github.com/tanyiqin/zack/znet"
)

//1 登录处理函数
func AccountLogin(request znet.IRequest) {
	msg := &pb.AccountLogin{}
	err := proto.Unmarshal(request.GetMsg().GetMsgData(), msg)
	if err != nil {
		fmt.Println("proto unmarshal error", err)
		return
	}
	fmt.Println("Rolid=", msg.RoleID, "Pwd=", string(msg.PassWord))
	msgReturn := &pb.AccountLoginResult{Result: 1}
	request.GetConn().SendMsg(1, msgReturn)
}
