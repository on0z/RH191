package main

import (
	"fmt"

	librh191 "github.com/on0z/RH191"
	"github.com/on0z/RH191/types"
)

func main() {
	r := librh191.NewRH191(types.On, types.ModeCool, 26)
	hex := r.GetHex()
	fmt.Println(hex)
	bin := r.GetBinary()
	fmt.Println(bin)

	adrsir := genADRSIR(bin)
	fmt.Println(adrsir)
}

/// HexからADRSIRコードを生成する
/// - Parameter hex: 送信するHex
/// - Returns: ADRSIRコード
func genADRSIR(binary string) string {
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
