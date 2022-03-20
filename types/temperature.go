package types

import "github.com/pkg/errors"

const (
	TEMPERATURE_Pos      = 0
	TEMPERATURE_Msk byte = 0b00001111
)

type Temperature int

func (t *Temperature) Validation() error {
	if !(16 <= *t && *t <= 31) {
		return errors.New("invalid temperature")
	}
	return nil
}

func (t *Temperature) GetFlag() byte {
	return byte(*t - 16)
}
