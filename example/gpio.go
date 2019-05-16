package main

import (
	"fmt"
	"os"

	"github.com/kshamko/gopherberry"
)

func main() {

	fmt.Println(0xF2200000, os.Getpagesize(), 0xF2200000%os.Getpagesize())

	r, _ := gopherberry.New()
	p17, _ := r.GetPin(11)
	err := p17.SetHigh()

	if err != nil {
		fmt.Println("[ERROR] cant sethight to pin 17(11)")
	}
}

//00 000 000 001 000 000 000 000 000 000 000
//1110 0000 0000 0000 0000 0000
