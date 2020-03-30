package main

import (
	"fmt"
	"os"

	"github.com/kshamko/gopherberry"
)

func main() {
	r, err := gopherberry.New(gopherberry.ARM2837)
	if err != nil {
		fmt.Println("[ERROR] can't init pi", err)
	}

	p17, _ := r.GetPin(12)
	_ = p17.ModeOutput()

	p40, _ := r.GetPin(40)
	p40.ModeInput()

	c, err := p40.DetectEdge(gopherberry.EdgeBoth)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for {
		edge := <-c

		if edge == gopherberry.EdgeHigh {
			fmt.Println("rising")
			p17.SetHigh()
		}

		if edge == gopherberry.EdgeLow {
			fmt.Println("falling")
			p17.SetLow()
		}
	}

}
