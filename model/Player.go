package model

import "github.com/tanyiqin/zack/znet"

type Player struct {
	Conn znet.IConnection
	Name string
}

func NewPlayer(conn znet.IConnection, name string) *Player{
	p := &Player{
		Conn: conn,
		Name: name,
	}
	return p
}

