package gopherberry

import (
	"testing"

	"gotest.tools/assert"
)

func TestOffsets(t *testing.T) {
	c := newChip2837()
	registers := c.getGPIORegisters()
	o := offsets(registers, registers["GPFSEL"][0]) //, c.getBasePeriphialsAddress())
	assert.Equal(t, o[registers["GPFSEL"][0]], 0)
	assert.Equal(t, o[registers["GPSET"][0]], 7)
}
