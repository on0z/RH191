package librh191

import (
	"fmt"
	"math/bits"

	"github.com/on0z/RH191/types"
)

type Active = types.Active
type Mode = types.Mode
type Temperature = types.Temperature
type Speed = types.Speed
type Direction = types.Direction
type Sound = types.Sound

type RH191API interface {
	SetActive(active Active)
	SetMode(mode Mode)
	SetTemperature(temperature Temperature) error
	SetSpped(speed Speed)
	SetDirection(direction Direction)
	SetSound(soundCnt Sound)
	GetBodyBytesSlice() []byte
	GetBytesSlice() []byte
	GetHex() string
	GetBinary() string
}

type rh191 struct {
	// 0~4バイト目
	// Default: 0x23, 0xCB, 0x26, 0x01, 0x00
	initialBytes []byte
	// 5バイト目
	// Default: 0x00
	ActiveReg byte
	// 6バイト目
	// Default: 0x00
	Mode1Reg byte
	// 7バイト目
	// Default: 0x00
	TemperatureReg byte
	// 8バイト目
	// Default: 0x30
	Mode2Reg byte
	// 9バイト目
	// Default: 0x00
	OtherConfigReg byte
	// 10~16バイト目
	// Default: 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00
	blangReg []byte
	// 17バイト目
	CheckReg byte
}

func NewRH191(active Active, mode Mode, temperature Temperature) RH191API {
	r := &rh191{
		initialBytes: []byte{0x23, 0xCB, 0x26, 0x01, 0x00},
		Mode2Reg:     0x30,
		blangReg:     []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
	}
	r.SetActive(active)
	r.SetMode(mode)
	r.SetTemperature(temperature)
	r.SetSound(types.SOUNT_COUNT1)
	return r
}

func (r *rh191) SetActive(active Active) {
	r.ActiveReg &= ^types.ACTIVE_Msk
	r.ActiveReg |= active.GetFlag() << types.ACTIVE_Pos
}

func (r *rh191) SetMode(mode Mode) {
	r.Mode1Reg &= ^types.MODE1_Msk
	r.Mode1Reg |= mode.GetFlag1() << types.MODE1_Pos
	r.Mode2Reg &= ^types.MODE2_Msk
	r.Mode2Reg |= mode.GetFlag2() << types.MODE2_Pos

	if mode == types.MODE_DRY {
		r.SetTemperature(24)
	}
}

func (r *rh191) SetTemperature(temperature Temperature) error {
	if err := temperature.Validation(); err != nil {
		return err
	}

	r.TemperatureReg &= ^types.TEMPERATURE_Msk
	r.TemperatureReg |= temperature.GetFlag() << types.TEMPERATURE_Pos

	return nil
}

func (r *rh191) SetSpped(speed Speed) {
	r.OtherConfigReg &= ^types.SPEED_Msk
	r.OtherConfigReg |= speed.GetFlag() << types.SPEED_Pos
}

func (r *rh191) SetDirection(direction Direction) {
	r.OtherConfigReg &= ^types.DIRECTION_Msk
	r.OtherConfigReg |= ^direction.GetFlag() << types.DIRECTION_Pos
}

func (r *rh191) SetSound(soundCnt Sound) {
	r.OtherConfigReg &= ^types.SOUND_Msk
	r.OtherConfigReg |= soundCnt.GetFlag() << types.SOUND_Pos
}

func (r *rh191) updateCheckReg() {
	bytes := r.GetBodyBytesSlice()

	var sum byte = 0
	for _, b := range bytes {
		sum += b
	}

	r.CheckReg = sum
}

func (r *rh191) GetBodyBytesSlice() []byte {
	bytes := []byte{}
	bytes = append(bytes, r.initialBytes...)
	bytes = append(bytes, r.ActiveReg)
	bytes = append(bytes, r.Mode1Reg)
	bytes = append(bytes, r.TemperatureReg)
	bytes = append(bytes, r.Mode2Reg)
	bytes = append(bytes, r.OtherConfigReg)
	bytes = append(bytes, r.blangReg...)
	return bytes
}

func (r *rh191) GetBytesSlice() []byte {
	bytes := r.GetBodyBytesSlice()
	bytes = append(bytes, r.CheckReg)
	return bytes
}

func (r *rh191) GetHex() string {
	r.updateCheckReg()
	s := r.GetBytesSlice()

	str := ""

	for _, b := range s {
		reved := bits.Reverse8(b)
		str += fmt.Sprintf("%02X", reved)
	}

	return str
}

func (r *rh191) GetBinary() string {
	r.updateCheckReg()
	s := r.GetBytesSlice()

	str := ""

	for _, b := range s {
		reved := bits.Reverse8(b)
		str += fmt.Sprintf("%08b", reved)
	}

	return str
}

func reverse(s []byte) []byte {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
