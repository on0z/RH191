package librh191

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/on0z/RH191/types"
)

func TestRH191GetHex(t *testing.T) {

	type input struct {
		active         types.Active
		mode           types.Mode
		temperature    types.Temperature
		isSetSpeed     bool
		speed          types.Speed
		isSetDirection bool
		direction      types.Direction
		isSetSound     bool
		sound          types.Sound
	}

	type expect struct {
		Res string
	}

	cases := []struct {
		name string // モード 温度 風速 風向
		input
		expect
	}{
		{
			name: "冷房 26 Auto Auto",
			input: input{
				active:      types.On,
				mode:        types.ModeCool,
				temperature: 26,
			},
			expect: expect{
				Res: "C4D36480000418506C0200000000000000B3C4D36480000418506C0200000000000000B3",
			},
		},
		{
			name: "冷房 27 Auto Auto",
			input: input{
				active:      types.On,
				mode:        types.ModeCool,
				temperature: 27,
				isSetSound:  true,
				sound:       types.SoundCount2,
			},
			expect: expect{
				Res: "C4D36480000418D06C010000000000000070C4D36480000418D06C010000000000000070",
			},
		},
		{
			name: "除湿  Auto Auto",
			input: input{
				active:      types.On,
				mode:        types.ModeDry,
				temperature: 24,
			},
			expect: expect{
				Res: "C4D36480000408104C0200000000000000FDC4D36480000408104C0200000000000000FD",
			},
		},
		{
			name: "暖房 20 Auto Auto",
			input: input{
				active:      types.On,
				mode:        types.ModeHeat,
				temperature: 20,
			},
			expect: expect{
				Res: "C4D36480000410200C02000000000000008DC4D36480000410200C02000000000000008D",
			},
		},
		{
			name: "Off (冷26AutoAuto)",
			input: input{
				active:      types.Off,
				mode:        types.ModeCool,
				temperature: 26,
			},
			expect: expect{
				Res: "C4D36480000018506C0200000000000000B5C4D36480000018506C0200000000000000B5",
			},
		},
		{
			name: "冷房 26 1 Auto",
			input: input{
				active:      types.On,
				mode:        types.ModeCool,
				temperature: 26,
				isSetSpeed:  true,
				speed:       types.SpeedWeak,
			},
			expect: expect{
				Res: "C4D36480000418506C820000000000000073C4D36480000418506C820000000000000073",
			},
		},
		{
			name: "冷房 26 2 Auto",
			input: input{
				active:      types.On,
				mode:        types.ModeCool,
				temperature: 26,
				isSetSpeed:  true,
				speed:       types.SpeedMiddle,
			},
			expect: expect{
				Res: "C4D36480000418506C4200000000000000F3C4D36480000418506C4200000000000000F3",
			},
		},
		{
			name: "冷房 26 3 Auto",
			input: input{
				active:      types.On,
				mode:        types.ModeCool,
				temperature: 26,
				isSetSpeed:  true,
				speed:       types.SpeedStrong,
			},
			expect: expect{
				Res: "C4D36480000418506CC2000000000000000BC4D36480000418506CC2000000000000000B",
			},
		},
		{
			name: "冷房 26 Auto 上",
			input: input{
				active:         types.On,
				mode:           types.ModeCool,
				temperature:    26,
				isSetDirection: true,
				direction:      types.DirectionUp,
			},
			expect: expect{
				Res: "C4D36480000418506C1200000000000000ABC4D36480000418506C1200000000000000AB",
			},
		},
		{
			name: "冷房 26 Auto 中上",
			input: input{
				active:         types.On,
				mode:           types.ModeCool,
				temperature:    26,
				isSetDirection: true,
				direction:      types.DirectionMiddleUp,
			},
			expect: expect{
				Res: "C4D36480000418506C0A00000000000000BBC4D36480000418506C0A00000000000000BB",
			},
		},
		{
			name: "冷房 26 Auto 中",
			input: input{
				active:         types.On,
				mode:           types.ModeCool,
				temperature:    26,
				isSetDirection: true,
				direction:      types.DirectionMiddle,
			},
			expect: expect{
				Res: "C4D36480000418506C1A00000000000000A7C4D36480000418506C1A00000000000000A7",
			},
		},
		{
			name: "冷房 26 Auto 中下",
			input: input{
				active:         types.On,
				mode:           types.ModeCool,
				temperature:    26,
				isSetDirection: true,
				direction:      types.DirectionMiddleDown,
			},
			expect: expect{
				Res: "C4D36480000418506C0600000000000000B7C4D36480000418506C0600000000000000B7",
			},
		},
		{
			name: "冷房 26 Auto 下",
			input: input{
				active:         types.On,
				mode:           types.ModeCool,
				temperature:    26,
				isSetDirection: true,
				direction:      types.DirectionDown,
			},
			expect: expect{
				Res: "C4D36480000418506C1600000000000000AFC4D36480000418506C1600000000000000AF",
			},
		},
		{
			name: "冷房 26 Auto スイング",
			input: input{
				active:         types.On,
				mode:           types.ModeCool,
				temperature:    26,
				isSetDirection: true,
				direction:      types.DirectionSwing,
			},
			expect: expect{
				Res: "C4D36480000418506C1E00000000000000A0C4D36480000418506C1E00000000000000A0",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := NewRH191(c.input.active, c.input.mode, c.input.temperature)
			if c.input.isSetSpeed {
				r.SetSpped(c.input.speed)
			}
			if c.input.isSetDirection {
				r.SetDirection(c.input.direction)
			}
			if c.input.isSetSound {
				r.SetSound(c.input.sound)
			}
			res := r.GetHex()
			assert.Equal(t, c.expect.Res, res+res)
		})
	}

}
