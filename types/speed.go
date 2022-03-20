package types

const (
	SPEED_Pos      = 0
	SPEED_Msk byte = 0b00000011
)

type Speed uint8

const (
	SpeedAuto Speed = iota
	SpeedWeak
	SpeedMiddle
	SpeedStrong
)

func (t *Speed) GetFlag() byte {
	return byte(*t)
}
