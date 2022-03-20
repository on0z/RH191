# RH191 for golang

[for swift](#for-swift)

三菱のエアコンのリモコン RH191 の信号の16進数文字列を生成するスクリプト

# How to use

see more: [cmd/sample/main.go](./cmd/sample/main.go)

ビットトレードワン赤外線送受信機 ADRSIR を使ってエアコンを操作する場合

```golang
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

func genADRSIRCmd(binary string) string {
    ...
}
```

# for swift

https://github.com/on0z/RH191/tree/release/swift

# 参考
リモコンの信号の中身を調べた時のメモ https://qiita.com/on0z/items/4d71ecdee7db3d44a8a9
