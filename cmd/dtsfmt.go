package main

import (
	"dtsfmt"
	"fmt"
)

func main() {
	var p dtsfmt.Prettier
	formatted, err := p.Prettify("corpus/linux/arch/arm64/boot/dts/intel/keembay-evm.dts")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(formatted))
}
