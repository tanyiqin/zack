package utils

const (
	IDRole = iota //角色ID
)

type idFactory struct {
	idList map[uint32]uint32
}


