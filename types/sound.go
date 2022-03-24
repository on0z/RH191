package types

const (
	SOUND_Pos      = 6
	SOUND_Msk byte = 0b11000000
)

type Sound uint8

const (
	SoundCount1 Sound = iota
	SoundCount2
)

func (t *Sound) GetFlag() byte {
	return byte(*t + 1)
}
