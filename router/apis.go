package router

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	log "github.com/tanyiqin/zack/logger"
	"github.com/tanyiqin/zack/model"
	"github.com/tanyiqin/zack/pb"
	"github.com/tanyiqin/zack/znet"
)

//1 登录处理函数
func CsAccountLogin(request znet.IRequest) {
	msg := &pb.CsAccountLogin{}
	err := proto.Unmarshal(request.GetMsg().GetMsgData(), msg)
	if err != nil {
		log.Error("proto unmarshal error", err)
		return
	}
	// 这边肯定去mongo数据库去验证 是否有改角色
	fmt.Println("Rolid=", msg.RoleID, "Pwd=", string(msg.PassWord))
	// 这里先假设搞一个角色
	p := model.NewPlayer(request.GetConn(), "你妈")
	request.GetConn().SetProperty("player", p)

	msgReturn := &pb.ScAccountLoginResult{Result: 1}
	request.GetConn().SendMsg(1, msgReturn)
}

//2 创建角色
func CsPlayerCreate(request znet.IRequest) {
	msg := &pb.CsPlayerCreate{}
	err := proto.Unmarshal(request.GetMsg().GetMsgData(), msg)
	if err != nil {
		log.Error("proto unmarshal error", err)
		return
	}

	// 这里要到数据库去创角
	player := model.NewPlayer(request.GetConn(), "你妈")
	request.GetConn().SetProperty("player", player)

	msgReturn := &pb.ScPlayerCreateResult{Result: 1}
	request.GetConn().SendMsg(2, msgReturn)
}

