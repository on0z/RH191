package types

const (
	DIRECTION_Pos      = 3
	DIRECTION_Msk byte = 0b00111000
)

type Direction uint8

const (
	DirectionAuto Direction = iota
	DirectionUp
	DirectionMiddleUp
	DirectionMiddle
	DirectionMiddleDown
	DirectionDown
	_
	DirectionSwing
)

func (t *Direction) GetFlag() byte {
	return byte(*t)
}
