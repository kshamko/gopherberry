package main

import (
	"fmt"

	"github.com/kshamko/gopherberry"
)

func main() {

	r, err := gopherberry.New(gopherberry.ARM2837)

	if err != nil {
		fmt.Println("[ERROR] can't init pi", err)
	}

	p17, _ := r.GetPin(11)
	err = p17.ModeOutput()
	if err != nil {
		fmt.Println("[ERROR] cant set mode to pin 17(11)", err)
	}

	err = p17.SetHigh()
	if err != nil {
		fmt.Println("[ERROR] cant sethigh to pin 17(11)", err)
	}
}

//00 000 000 001 000 000 000 000 000 000 000
//1110 0000 0000 0000 0000 0000
