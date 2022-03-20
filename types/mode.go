package types

const (
	MODE1_Pos      = 3
	MODE1_Msk byte = 0b00011000
	MODE2_Pos      = 1
	MODE2_Msk byte = 0b00000110
)

type Mode uint8

const (
	MODE_HEAT Mode = iota
	MODE_DRY
	MODE_COOL
)

func (t *Mode) GetFlag1() byte {
	return byte(*t + 1)
}

func (t *Mode) GetFlag2() byte {
	switch *t {
	case MODE_HEAT:
		return 0
	case MODE_DRY:
		return 1
	case MODE_COOL:
		return 3
	default:
		return 0
	}
}
