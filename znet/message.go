package znet

var HeadLen uint32 = 8 // Message 头部字段所占字节

// 消息组成 4字节消息体长度 4字节消息ID x字节proto编码的消息体

type IMessage interface {
	GetMsgID() uint32
	GetDataLen() uint32
	GetMsgData() []byte
	SetMsgData([]byte)
}

type Message struct {
	DataLen	uint32
	MsgID uint32
	Data []byte
}

func NewMessage(msgID uint32, data []byte) IMessage{
	s := &Message{
		Data: data,
		DataLen: uint32(len(data)),
		MsgID: msgID,
	}
	return s
}

func (m *Message)GetMsgID() uint32 {
	return m.MsgID
}
func (m *Message)GetDataLen() uint32 {
	return m.DataLen
}
func (m *Message)GetMsgData() []byte {
	return m.Data
}

func (m *Message)SetMsgData(data []byte) {
	m.Data = data
}

