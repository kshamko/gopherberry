package main

import (
	"fmt"
	//"github.com/kshamko/gopherberry"
)

//See GPFSELn spec for details
const (
	MODE_INPUT  = 0
	MODE_OUTPUT = 1

	MODE_ALT0 = 4
	MODE_ALT1 = 5
	MODE_ALT2 = 6
	MODE_ALT3 = 7
	MODE_ALT4 = 3
	MODE_ALT5 = 2
)

func main() {

	fmt.Printf("00%b\n", MODE_INPUT)
	fmt.Printf("00%b\n", MODE_OUTPUT)
	fmt.Printf("%b\n", MODE_ALT0)
	fmt.Printf("%b\n", MODE_ALT1)
	fmt.Printf("%b\n", MODE_ALT2)
	fmt.Printf("%b\n", MODE_ALT3)
	fmt.Printf("0%b\n", MODE_ALT4)
	fmt.Printf("0%b\n", MODE_ALT5)

	//baseCommand := uint8(0)

	/*chip, err := gopherberry.NewChip()
	fmt.Printf("%+v, %v", chip, err)*/

	pin := 17
	//fselReg := uint8(pin) / 10
	shift := (uint8(pin) % 10) * 3
	//f := uint32(0)

	//const pinMask = 7 // 111 - pinmode is 3 bits

	fmt.Println("Shift: ", shift)
	fmt.Printf("%b\n", MODE_OUTPUT<<shift)
}

//00 000 000 001 000 000 000 000 000 000 000
//1110 0000 0000 0000 0000 0000
