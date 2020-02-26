package gopherberry

import (
	"fmt"
	"testing"
)

func TestStopClock(t *testing.T) {

	pi, _ := New(ARM2837)
	pin, _ := pi.GetPin(12)

	pin.StopClock()
	const PASSWORD = 0x5A000000
	const enab = 1 << 4
	mash := 0
	src := 1 << 0
	command := PASSWORD | mash | src | enab
	fmt.Printf("%b\n", command)
}
