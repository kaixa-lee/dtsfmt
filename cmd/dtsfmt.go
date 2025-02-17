package main

import (
	"dtsfmt"
	"flag"
	"fmt"
)

func main() {
	useDebug := flag.Bool("debug", false, "")
	integerCellSize := flag.Int("icz", 4, "")

	flag.Parse()

	opts := []dtsfmt.Option{
		dtsfmt.WithDebug(*useDebug),
		dtsfmt.WithIntegerCellSize(*integerCellSize),
	}

	var p dtsfmt.Prettier
	for _, f := range opts {
		f(&p)
	}

	formatted, err := p.Prettify("corpus/linux/arch/arm64/boot/dts/intel/keembay-evm.dts")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(formatted))
}
