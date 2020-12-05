package router

import (
	"github.com/golang/protobuf/proto"
	mdb "github.com/tanyiqin/zack/db"
	log "github.com/tanyiqin/zack/logger"
	"github.com/tanyiqin/zack/model"
	"github.com/tanyiqin/zack/pb"
	"github.com/tanyiqin/zack/znet"
	"go.mongodb.org/mongo-driver/bson"
)

//1 登录处理函数
func CsAccountLogin(request znet.IRequest) {
	msg := &pb.CsAccountLogin{}
	msgReturn := &pb.ScAccountLoginResult{}
	err := proto.Unmarshal(request.GetMsg().GetMsgData(), msg)
	if err != nil {
		log.Error("proto unmarshal error", err)
		return
	}
	// 这边去mongo数据库去验证 是否有该角色
	findResult := mdb.DB.FindOne("g_role", bson.M{"_id":msg.RoleID})
	var p model.Player
	err = findResult.Decode(&p)
	if err != nil {
		msgReturn.Result = 2
		request.GetConn().SendMsg(1, msgReturn)
		return
	}
	request.GetConn().SetProperty("player", p)
	msgReturn.Result = 1
	request.GetConn().SendMsg(1, msgReturn)
}

//2 创建角色
func CsPlayerCreate(request znet.IRequest) {
	msg := &pb.CsPlayerCreate{}
	msgReturn := &pb.ScPlayerCreateResult{}
	err := proto.Unmarshal(request.GetMsg().GetMsgData(), msg)
	if err != nil {
		log.Error("proto unmarshal error", err)
		return
	}

	// 这里要到数据库去创角
	player := model.NewPlayer(request.GetConn(), msg.Name)
	data, err := bson.Marshal(player)
	if err != nil {
		log.Error("CreatePlayerMarshal error", err)
		msgReturn.Result = 2
		return
	}
	_, err = mdb.DB.InsertOne("g_role", data)
	if err != nil {
		log.Error("mongo createRole err", err)
		msgReturn.Result = 2
		return
	}
	request.GetConn().SetProperty("player", player)

	msgReturn.Result = 1
	request.GetConn().SendMsg(2, msgReturn)
}

