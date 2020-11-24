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
}

func NewServer() IServer {
	s := &Server{
		IP: "0.0.0.0",
		Port: 7777,
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

// 启动服务器
func (s *Server)Start() {
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.IP, s.Port))

	if err != nil {
		fmt.Println("listen error = ", err)
		return
	}

	// 唯一ID  这个后续考虑改成用mongo数据库的自增键
	var cid uint32 = 0
	for {
		// 这里后续改成连接池处理
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