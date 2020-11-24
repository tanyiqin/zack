package znet

import "sync"

// 管理当前server下的所有链接
type IConnMgr interface {
	// 添加conn
	AddConn(connection IConnection)
	// 删除conn
	Remove(connID uint32)
	// 停止conn
	Stop()
}

type ConnMgr struct {
	connections map[uint32]IConnection
	connLock sync.RWMutex
}

func NewConnManager() IConnMgr{
	return &ConnMgr{
		connections: make(map[uint32]IConnection),
	}
}

// 添加conn
func (c *ConnMgr)AddConn(connection IConnection) {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	c.connections[connection.GetConnID()] = connection
}
// 删除conn
func (c *ConnMgr)Remove(connID uint32) {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	delete(c.connections, connID)
}
// 停止所有conn
func (c *ConnMgr)Stop() {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	for _, conn := range c.connections {
		conn.Stop()
	}
	// 重新指向新的空间 原先的数据由垃圾回收处理
	c.connections = make(map[uint32]IConnection)
}
