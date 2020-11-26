package znet

import (
	"fmt"
	"net"
)

// 一切的起点
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

// 添加Router
func (s *Server)AddRouter(msgID uint32, handleFunc HandleFunc) {
	s.MsgHandler.AddRouter(msgID, handleFunc)
}

// 启动服务器
func (s *Server)Start() {
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.IP, s.Port))
	defer listen.Close()

	if err != nil {
		fmt.Println("listen error = ", err)
		return
	}

	// 唯一ID  这个后续考虑改成用mongo数据库的自增键
	var cid uint32 = 0
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Accept err = ", err)
			continue
		}

		dealConn := NewConnection(s, conn, cid)
		cid++

		go dealConn.Start()
	}
}

// 关闭服务器
func (s *Server)Stop() {

	// 移除所有管理器
	s.ConnMgr.Stop()
}