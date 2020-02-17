package gopherberry

import (
	"testing"
	"fmt"
)

func TestOffsets(t *testing.T) {
	/*c := newChip2837()
	registers := c.getGPIORegisters()
	o := offsets(registers, registers["GPFSEL"][0]) //, c.getBasePeriphialsAddress())
	assert.Equal(t, o[registers["GPFSEL"][0]], 0)
	assert.Equal(t, o[registers["GPSET"][0]], 7)*/
}

func TestStopPWM(t *testing.T) {

	r, _ := New(ARM2837)
	r.StopPWM()

	fmt.Printf("%b !!%b!!\n", (1 << 8 | 1), 0b00000000000000000000000100000001 &^ (1<<8 | 1))
}
