package znet

import (
	log "github.com/tanyiqin/zack/logger"
)

// 消息处理模块 目前为同步执行

// 消息处理函数
type HandleFunc func(request IRequest)

type IMsgHandler interface {
	// 添加路由消息
	AddRouter(msgID uint32, handleFunc HandleFunc)
	// 执行路由消息
	DoMsgRouter(request IRequest)
}

type MsgHandler struct {
	RouterMap	map[uint32]HandleFunc
}

func NewMsgManager() IMsgHandler{
	s := &MsgHandler{
		RouterMap: make(map[uint32]HandleFunc),
	}
	return s
}

// 添加路由消息
func (mh *MsgHandler)AddRouter(msgID uint32, handleFunc HandleFunc) {
	if _, ok := mh.RouterMap[msgID]; ok {
		log.Debug("warning dump router")
	}
	mh.RouterMap[msgID] = handleFunc
}

// 执行路由消息
func (mh *MsgHandler) DoMsgRouter(request IRequest) {
	handFunc, ok := mh.RouterMap[request.GetMsg().GetMsgID()]
	if !ok {
		log.Error("Wrong msgID ", request.GetMsg().GetMsgID())
		return
	}

	handFunc(request)
}