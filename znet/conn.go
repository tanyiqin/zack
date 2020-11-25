package znet

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
	"net"
	"sync"
)

type IConnection interface {
	// 执行链接逻辑
	Start()
	// 关闭连接
	Stop()
	// 启动读逻辑
	StartReader()
	// 启动写逻辑
	StartWriter()
	// 获取唯一的ConnID
	GetConnID()uint32
	// 发送数据给客户端
	SendMsg(msgID uint32, message proto.Message)
	SendRawMsg(msgID uint32, data []byte) error

	// 设置额外的属性字段
	SetProperty(key string, value interface{})
	// 获取属性字段
	GetProperty(key string) (interface{}, error)
	// 移除属性字段
	RemoveProperty(key string)
}

type Connection struct {
	Server IServer
	// 链接套接字
	Conn net.Conn
	// 关闭chan
	ExitChan chan bool
	// 消息传输读，写 协程通讯chan
	MsgChan chan []byte
	// 链接ID 唯一
	ConnID uint32
	// 链接是否已经进行过stop
	isClosed bool
	// 属性锁
	propertyLock sync.RWMutex
	// 额外属性
	property map[string]interface{}
}

func NewConnection(server IServer, conn net.Conn, connID uint32) IConnection {
	s := &Connection{
		Server: server,
		Conn: conn,
		ExitChan: make(chan bool),
		MsgChan: make(chan []byte),
		ConnID: connID,
		isClosed: false,
		property: make(map[string]interface{}),
	}

	s.Server.GetConnMgr().AddConn(s)
	return s
}

// 开启链接读写逻辑
func (c *Connection)Start() {
	go c.StartReader()
	go c.StartWriter()
}

// 关闭连接
func (c *Connection)Stop() {

	// 防止多次执行Stop函数
	if c.isClosed {
		return
	}
	c.isClosed = true

	// 执行Hook函数
	c.Server.CallOnConnStop(c)

	// 关闭写进程
	c.ExitChan<-true

	// 关闭套接字
	fmt.Println("conn close")
	c.Conn.Close()

	c.Server.GetConnMgr().Stop()

	// 关闭管道
	close(c.ExitChan)
	close(c.MsgChan)
}

func (c *Connection)GetConnID() uint32{
	return c.ConnID
}

// 启动读逻辑
func (c *Connection)StartReader() {
	defer c.Stop()

	fmt.Println("Conn Reader begin connid=", c.GetConnID())

	for {
		msg, err := ReadFromConn(c.Conn)
		if err != nil {
			return
		}

		// 当前处理 逻辑 使用同步数据的方式 保证数据顺序的有序

		c.Server.GetMsgHandler().DoMsgRouter(&Request{Msg: msg, Conn: c})
	}
}

func ReadFromConn(Conn net.Conn) (IMessage, error){
	headData := make([]byte, HeadLen)

	// 读取头部8字节内容 消息长度+消息ID
	_, err := io.ReadFull(Conn, headData)
	if err != nil {
		fmt.Println("read msg head err", err)
		return nil, err
	}

	msg, err := UnPack(headData)
	if err != nil {
		fmt.Println("unpack err", err)
		return nil, err
	}

	// 根据头部提供的长度 读取data信息
	var data []byte
	if msg.GetDataLen() > 0 {
		data = make([]byte, msg.GetDataLen())
		if _, err := io.ReadFull(Conn, data); err != nil {
			fmt.Println("read msg err", err)
			return nil, err
		}
	}
	msg.SetMsgData(data)
	return msg, nil
}

// 启动写逻辑
// ！！这里如果写失败 把连接给中断掉
func (c *Connection)StartWriter() {
	defer c.Stop()
	fmt.Println("Conn Writer begin connid=", c.GetConnID())
	for {
		select {
		case data := <-c.MsgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("write err", err)
				return
			}
		case <-c.ExitChan:
			return
		}
	}
}

// 发送数据给客户端
func (c *Connection)SendMsg(msgID uint32, message proto.Message) {
	data, err := proto.Marshal(message)
	if err != nil {
		fmt.Println("proto marshal err,", err)
		return
	}

	if err := c.SendRawMsg(msgID, data); err != nil {
		return
	}
}
func (c *Connection)SendRawMsg(msgID uint32, data []byte) error {
	if c.isClosed {
		fmt.Println("conn already closed")
		return errors.New("conn already closed")
	}

	msg, err := Pack(NewMessage(msgID, data))

	if err != nil {
		fmt.Println("pack err = ", err)
		return errors.New("pack error msg")
	}

	c.MsgChan <- msg
	return nil
}

// 设置额外的属性字段
func (c *Connection)SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}
// 获取属性字段
func (c *Connection)GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property")
	}
}
// 移除属性字段
func (c *Connection)RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}