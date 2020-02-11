package gopherberry

import (
	"testing"

	"gotest.tools/assert"
)

func TestAddressOffset(t *testing.T) {

}

func TestAddrBus2Physt(t *testing.T) {
	c := newChip2837()
	assert.Equal(t, c.addrBus2Phys(0x7E200000), uint64(0x3F200000))
}

func TestPwmCtl(t *testing.T) {
	c := newChip2837()

	cfg1 := PWMChannelConfig{
		MSEnable:    1,
		ChanEnabled: 1,
	}

	cfg2 := PWMChannelConfig{}

	address, addressType, command := c.pwmCtl(cfg1, cfg2)
	assert.Equal(t, address, uint64(0x7E20C000))
	assert.Equal(t, command, 129) //129 = 0b0000000000000000000000010000001
	assert.Equal(t, addressType, addrBus)
}
