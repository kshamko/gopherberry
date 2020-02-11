package gopherberry

import (
	"testing"

	"gotest.tools/assert"
)

func TestMmapParameters(t *testing.T) {

	addresses := []uint64{
		0x7E20C020,
		0x7E20C000,
		0x7E20C004,
		0x7E20C008,
		0x7E20C010,
		0x7E20C018,
		0x7E20C024,
		0x7E20C014,
	}

	base, lenght, offsets := mmapParameters(addresses)

	assert.Equal(t, base, int64(0x7E20C000))
	assert.Equal(t, lenght, 8)
	assert.Equal(t, offsets[0x7E20C000], 0)
	assert.Equal(t, offsets[0x7E20C010], 3)
}
