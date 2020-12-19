package znet

import (
	"fmt"
	"github.com/tanyiqin/zack/logger"
	"net"
)

type IServer interface {
	// 启动服务器
	Start()
	// 停止服务器
	Stop()
	// 获取链接管理器
	GetConnMgr() IConnMgr
	// 获取消息处理器
	GetMsgHandler() IMsgHandler
	// 添加router
	AddRouter(msgID uint32, handleFunc HandleFunc)
	// Add链接中断时Hook
	AddConnStopFunc(func(IConnection))
	// 执行hook函数
	CallOnConnStop(connection IConnection)
	// Add链接开始时Hook
	AddConnStartFunc(func(connection IConnection))
	//
	CallOnConnStart(connection IConnection)
	// Server开始结束执行函数
	AddServerStartFunc(func())
	AddServerStopFunc(func())
}

type Server struct {
	// 监听ip地址
	IP string
	// 绑定端口
	Port uint32
	// 连接管理器
	ConnMgr IConnMgr
	// 消息管理器
	MsgHandler IMsgHandler
	// 连接中断Hook
	ConnStopFunc func(IConnection)
	//
	ConnStartFunc func(IConnection)
	ServerStartFunc func()
	ServerStopFunc func()
}

func NewServer() IServer {
	s := &Server{
		IP: "0.0.0.0",
		Port: 7777,
		ConnMgr: NewConnManager(),
		MsgHandler: NewMsgManager(),
	}
	return s
}

func (s *Server)GetConnMgr() IConnMgr {
	return s.ConnMgr
}

// 获取消息处理器
func (s *Server) GetMsgHandler() IMsgHandler {
	return s.MsgHandler
}

func (s *Server)AddConnStopFunc(connFunc func(IConnection)) {
	s.ConnStopFunc = connFunc
}

func (s *Server)CallOnConnStop(connection IConnection) {
	if s.ConnStopFunc != nil {
		s.ConnStopFunc(connection)
	}
}

// Add链接开始时Hook
func (s *Server) AddConnStartFunc(connFunc func(connection IConnection)) {
	s.ConnStartFunc = connFunc
}
//
func (s *Server) CallOnConnStart(connection IConnection) {
	if s.ConnStartFunc != nil {
		s.ConnStartFunc(connection)
	}
}

// Server开始结束执行函数
func (s *Server)AddServerStartFunc(f func()) {
	s.ServerStartFunc = f
}
func (s *Server)AddServerStopFunc(f func()) {
	s.ServerStopFunc = f
}

// 添加Router
func (s *Server)AddRouter(msgID uint32, handleFunc HandleFunc) {
	s.MsgHandler.AddRouter(msgID, handleFunc)
}

// 启动服务器
func (s *Server)Start() {
	if s.ServerStartFunc != nil {
		s.ServerStartFunc()
	}

	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.IP, s.Port))
	defer listen.Close()

	if err != nil {
		log.Fatal("Listen error", err)
		return
	}

	// 唯一ID  这个后续考虑改成用mongo数据库的自增键
	var cid uint32 = 0
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Error("Accept error", err)
			continue
		}

		dealConn := NewConnection(s, conn, cid)
		cid++

		go dealConn.Start()
	}
}

// 关闭服务器
func (s *Server)Stop() {
	if s.ServerStopFunc != nil {
		s.ServerStopFunc()
	}
	// 移除所有管理器
	s.ConnMgr.Stop()
}