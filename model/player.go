package model

import (
	"github.com/tanyiqin/zack/znet"
	"math/rand"
)

type Player struct {
	Conn znet.IConnection `bson:"-"`
	RoleID uint32 `bson:"_id"`
	Name string
}

func NewPlayer(conn znet.IConnection, name string) *Player{
	p := &Player{
		Conn: conn,
		Name: name,
		RoleID: rand.Uint32(),
	}
	return p
}

