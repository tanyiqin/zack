package znet

// 将conn 于 msg封装到一起
type IRequest interface {
	// 获取Data
	GetMsg() IMessage
	// 获取链接
	GetConn() IConnection
}

type Request struct {
	Conn IConnection
	Msg IMessage
}

// 获取Data
func (r *Request)GetMsg() IMessage {
	return r.Msg
}
// 获取链接
func (r *Request)GetConn() IConnection {
	return r.Conn
}