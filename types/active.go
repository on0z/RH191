package types

const (
	ACTIVE_Pos      = 5
	ACTIVE_Msk byte = 0b00100000
)

type Active uint8

const (
	Off Active = iota
	On
)

func (t *Active) GetFlag() byte {
	return byte(*t)
}
