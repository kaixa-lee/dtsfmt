package main

import "fmt"

func main() {
	var p Prettier
	formatted, err := p.Prettify("test/corpus/linux/arch/arm64/boot/dts/intel/keembay-evm.dts")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(formatted))
}
