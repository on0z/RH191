package types

const (
	MODE1_Pos          = 3
	MODE1_Msk     byte = 0b00011000
	MODE2_Pos          = 1
	MODE2_Msk     byte = 0b00000110
	MODE2_Default byte = 0x30
)

type Mode uint8

const (
	ModeHeat Mode = iota
	ModeDry
	ModeCool
)

func (t *Mode) GetFlag1() byte {
	return byte(*t + 1)
}

func (t *Mode) GetFlag2() byte {
	switch *t {
	case ModeHeat:
		return 0
	case ModeDry:
		return 1
	case ModeCool:
		return 3
	default:
		return 0
	}
}
