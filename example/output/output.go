package main

import (
	"fmt"
	"time"

	"github.com/kshamko/gopherberry"
)

func main() {

	r, err := gopherberry.New(gopherberry.ARM2837)

	if err != nil {
		fmt.Println("[ERROR] can't init pi", err)
	}

	p17, err := r.GetPin(12)
	if err != nil {
		fmt.Println("cant fetch pin", err)
		return
	}
	err = p17.ModeOutput()
	if err != nil {
		fmt.Println("[ERROR] cant set mode to pin 17(11)", err)
	}
	p17.SetLow()
	
	err = p17.SetHigh()
	if err != nil {
		fmt.Println("[ERROR] cant set high to pin 17(11)", err)
	}

	p26, _ := r.GetPin(37)
	p26.ModeOutput()
	p26.SetLow()

	l, _ := p17.Level()
	fmt.Println("[INFO] pin17 level:", l)

	for i := 0; i < 5; i++ {
		p26.SetHigh()
		time.Sleep(500 * time.Millisecond)
		p26.SetLow()
		time.Sleep(500 * time.Millisecond)
	}

	l, _ = p17.Level()
	fmt.Println("[INFO] pin17 level:", l)

	l, _ = p26.Level()
	fmt.Println("[INFO] pin26 level:", l)

	time.Sleep(2000 * time.Millisecond)
	p17.SetLow()

}

//00 000 000 001 000 000 000 000 000 000 000
//1110 0000 0000 0000 0000 0000
