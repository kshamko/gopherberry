package main

import (
	"fmt"

	"github.com/kshamko/gopherberry"
)

//
// Trigger it with "$ date | sudo tee /dev/kmsg"
//
func main() {

	defer fmt.Println("Stopped")

	ep, _ := gopherberry.NewEpoll("/sys/class/gpio/gpio21/value")

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

			//if x == 3 {
			//	return
			//}
		}
	}
}
