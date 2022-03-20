package main

import (
	"log"

	libadrsir "github.com/on0z/libadrsir-go"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	host "periph.io/x/host/v3"

	librh191 "github.com/on0z/RH191"
	"github.com/on0z/RH191/types"
)

func main() {

	// setup periph.io host
	_, err := host.Init()
	if err != nil {
		log.Fatalf("failed to initialize periph: %v", err)
	}

	// Use i2creg I²C bus registry to find the first available I²C bus.
	b, err := i2creg.Open("")
	if err != nil {
		log.Fatal(err)
	}
	defer b.Close()

	// Dev is a valid conn.Conn.
	d := &i2c.Dev{Addr: uint16(libadrsir.ADDR), Bus: b}

	// Setup ADRSIR
	adrsir := libadrsir.NewADRSIR(d)

	// Configure RH191
	r := librh191.NewRH191(types.On, types.ModeHeat, 21)
	r.SetDirection(types.DirectionDown)
	// Generate Binary
	bin := r.GetBinary()

	// Convert binary to ADRSIR command
	adrsirCmd := genADRSIRCmd(bin)

	// Send command
	adrsir.Send(adrsirCmd)
}

/// binaryからADRSIRコードを生成する
/// - Parameter hex: 送信するHex
/// - Returns: ADRSIRコード
func genADRSIRCmd(binary string) string {
	/// 始端
	adrsir := "83004300"

	/// 1と0をADRSIR表記にする
	for _, bit := range binary {
		if bit == '1' {
			adrsir += "11003200"
		} else if bit == '0' {
			adrsir += "11001100"
		}
	}

	/// 終端
	adrsir += "1300EE01"

	/// 2回繰り返して返却する
	return adrsir + adrsir
}
