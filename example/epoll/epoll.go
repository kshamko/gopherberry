package main

import (
	"fmt"
	"os"

	"github.com/kshamko/gopherberry"
)

//
// Trigger it with "$ date | sudo tee /dev/kmsg"
//
func main() {

	defer fmt.Println("Stopped")

	ep, err := gopherberry.NewEpoll("/dev/kmsg")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("started")
	x := 0
	for {
		c := ep.Wait()
		select {
		case _, ok := <-c:

			if !ok {
				fmt.Println("closed")
				return
			}
			x++
			fmt.Println("changed", x, "times")

			if x == 3 {
				return
			}
		}
	}
}
