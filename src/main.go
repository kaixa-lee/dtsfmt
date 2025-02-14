package main

import "fmt"

func main() {
	var p Prettier
	formatted, err := p.Prettify("test.dts")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(formatted))
}
