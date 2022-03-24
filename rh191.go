package librh191

import (
	"fmt"
	"math/bits"

	"github.com/on0z/RH191/types"
)

type RH191API interface {
	GetHex(config types.CommandConfig) (string, error)
	GetBinary(config types.CommandConfig) (string, error)
}

type rh191 struct{}

type rh191Command struct {
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
	blankReg []byte
	// 17バイト目
	CheckReg byte
}

func NewRH191() RH191API {
	return &rh191{}
}

func (_ *rh191) setActive(cmd *rh191Command, active types.Active) *rh191Command {
	cmd.ActiveReg &= ^types.ACTIVE_Msk
	cmd.ActiveReg |= active.GetFlag() << types.ACTIVE_Pos

	return cmd
}

func (r *rh191) setMode(cmd *rh191Command, mode types.Mode) *rh191Command {
	cmd.Mode1Reg &= ^types.MODE1_Msk
	cmd.Mode1Reg |= mode.GetFlag1() << types.MODE1_Pos
	cmd.Mode2Reg &= ^types.MODE2_Msk
	cmd.Mode2Reg |= mode.GetFlag2() << types.MODE2_Pos

	if mode == types.ModeDry {
		r.setTemperature(cmd, 24)
	}

	return cmd
}

func (_ *rh191) setTemperature(cmd *rh191Command, temperature types.Temperature) (*rh191Command, error) {
	if err := temperature.Validation(); err != nil {
		return nil, err
	}

	cmd.TemperatureReg &= ^types.TEMPERATURE_Msk
	cmd.TemperatureReg |= temperature.GetFlag() << types.TEMPERATURE_Pos

	return cmd, nil
}

func (_ *rh191) setSpeed(cmd *rh191Command, speed types.Speed) *rh191Command {
	cmd.OtherConfigReg &= ^types.SPEED_Msk
	cmd.OtherConfigReg |= speed.GetFlag() << types.SPEED_Pos

	return cmd
}

func (_ *rh191) setDirection(cmd *rh191Command, direction types.Direction) *rh191Command {
	cmd.OtherConfigReg &= ^types.DIRECTION_Msk
	cmd.OtherConfigReg |= direction.GetFlag() << types.DIRECTION_Pos

	return cmd
}

func (_ *rh191) setSound(cmd *rh191Command, soundCnt types.Sound) *rh191Command {
	cmd.OtherConfigReg &= ^types.SOUND_Msk
	cmd.OtherConfigReg |= soundCnt.GetFlag() << types.SOUND_Pos

	return cmd
}

func (_ *rh191) updateCheckReg(cmd *rh191Command, bytes []byte) *rh191Command {
	var sum byte = 0
	for _, b := range bytes {
		sum += b
	}

	cmd.CheckReg = sum

	return cmd
}

func (_ *rh191) getBodyBytesSlice(cmd *rh191Command) []byte {
	bytes := []byte{}
	bytes = append(bytes, cmd.initialBytes...)
	bytes = append(bytes, cmd.ActiveReg)
	bytes = append(bytes, cmd.Mode1Reg)
	bytes = append(bytes, cmd.TemperatureReg)
	bytes = append(bytes, cmd.Mode2Reg)
	bytes = append(bytes, cmd.OtherConfigReg)
	bytes = append(bytes, cmd.blankReg...)
	return bytes
}

func (r *rh191) assembleBytesSlice(config types.CommandConfig) ([]byte, error) {
	cmd := &rh191Command{
		initialBytes: []byte{0x23, 0xCB, 0x26, 0x01, 0x00},
		Mode2Reg:     types.MODE2_Default,
		blankReg:     []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
	}

	r.setActive(cmd, config.Active)
	_, err := r.setTemperature(cmd, config.Temperature)
	if err != nil {
		return nil, err
	}
	r.setMode(cmd, config.Mode)
	r.setSpeed(cmd, types.Speed(config.Speed))
	r.setDirection(cmd, types.Direction(config.Direction))
	r.setSound(cmd, config.Sound)

	bytes := r.getBodyBytesSlice(cmd)
	r.updateCheckReg(cmd, bytes)
	bytes = append(bytes, cmd.CheckReg)
	return bytes, nil
}

func (r *rh191) GetHex(config types.CommandConfig) (string, error) {
	s, err := r.assembleBytesSlice(config)
	if err != nil {
		return "", err
	}

	str := ""

	for _, b := range s {
		reved := bits.Reverse8(b)
		str += fmt.Sprintf("%02X", reved)
	}

	return str, nil
}

func (r *rh191) GetBinary(config types.CommandConfig) (string, error) {
	s, err := r.assembleBytesSlice(config)
	if err != nil {
		return "", err
	}

	str := ""

	for _, b := range s {
		reved := bits.Reverse8(b)
		str += fmt.Sprintf("%08b", reved)
	}

	return str, nil
}

func reverse(s []byte) []byte {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
