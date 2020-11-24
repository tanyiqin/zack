package znet

import (
	"bytes"
	"encoding/binary"
)

// 用于将msg中的头部 ID 长度字段转换为-> []byte
func Pack(msg IMessage) ([]byte, error) {
	dataBuff := bytes.NewBuffer([]byte{})

	if err := binary.Write(dataBuff, binary.BigEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}

	if err := binary.Write(dataBuff, binary.BigEndian, msg.GetMsgID()); err != nil {
		return nil, err
	}

	if err := binary.Write(dataBuff, binary.BigEndian, msg.GetMsgData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

// 这里只解压msg固定的头部查看具体有多少的信息
func UnPack(binaryData []byte) (IMessage, error) {
	dataBuff := bytes.NewReader(binaryData)

	msg := &Message{}

	if err := binary.Read(dataBuff, binary.BigEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	if err := binary.Read(dataBuff, binary.BigEndian, &msg.MsgID); err != nil {
		return nil, err
	}


	return msg, nil
}
