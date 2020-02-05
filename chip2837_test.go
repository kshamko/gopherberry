package gopherberry

import (
	"fmt"
	"os"
	"testing"

	"gotest.tools/assert"
)

func TestAddressOffset(t *testing.T) {

}

func TestAddrBus2Physt(t *testing.T) {

	fmt.Println(os.Getpagesize())
	c := newChip2837()

	assert.Equal(t, c.addrBus2Phys(0x7E200000), uint64(0x3F200000))
}
