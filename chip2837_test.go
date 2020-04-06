package gopherberry

import (
	"testing"

	"gotest.tools/assert"
)

func TestGetPinBCM(t *testing.T) {
	testCases := []struct {
		name string
		in   int
		out  int
	}{
		{"success", 12, 18},
		{"error", 100, NoBCMNum},
	}

	c := new2837()
	for _, tt := range testCases {

		bcm := c.getPinBCM(tt.in)
		if bcm != tt.out {
			t.Errorf("Test case %s failed. Expected: %v. Got: %v\n", tt.name, tt.out, bcm)
		}
	}
}

func TestAddrBus2Physt(t *testing.T) {
	c := new2837()
	assert.Equal(t, c.addrBus2Phys(0x7E200000), uint64(0x3F200000))
}

func TestGpgsel(t *testing.T) {
	testCases := []struct {
		name      string
		bcm       int
		mode      PinMode
		address   uint64
		operation int
	}{
		{"p12 - alt5", 12, PinModeALT5, 0x7E200004, 128},
	}
	c := new2837()
	for _, tt := range testCases {

		addr, _, op := c.gpgsel(tt.bcm, tt.mode)
		if addr != tt.address || op != tt.operation {
			t.Errorf("Test case %s failed. Expected: %v, %b. Got: %v, %b\n", tt.name, tt.address, tt.operation, addr, op)
		}
	}
}

func TestGpset(t *testing.T) {

	testCases := []struct {
		name      string
		bcm       int
		address   uint64
		operation int
	}{
		{"p12", 12, 0x7E20001C, 4096},
		{"p34", 34, 0x7E200020, 4},
	}

	c := new2837()
	for _, tt := range testCases {

		addr, _, op := c.gpset(tt.bcm)
		if addr != tt.address || op != tt.operation {
			t.Errorf("Test case %s failed. Expected: %v, %b. Got: %v, %b\n", tt.name, tt.address, tt.operation, addr, op)
		}
	}
}

func TestGpclr(t *testing.T) {

	testCases := []struct {
		name      string
		bcm       int
		address   uint64
		operation int
	}{
		{"p12", 12, 0x7E200028, 4096},
		{"p34", 34, 0x7E20002C, 4},
	}

	c := new2837()
	for _, tt := range testCases {

		addr, _, op := c.gpclr(tt.bcm)
		if addr != tt.address || op != tt.operation {
			t.Errorf("Test case %s failed. Expected: %v, %b. Got: %v, %b\n", tt.name, tt.address, tt.operation, addr, op)
		}
	}
}

func TestGplev(t *testing.T) {

	testCases := []struct {
		name      string
		bcm       int
		address   uint64
		operation int
	}{
		{"p12", 12, 0x7E200034, 4096},
		{"p34", 34, 0x7E200038, 4},
	}

	c := new2837()
	for _, tt := range testCases {

		addr, _, op := c.gplev(tt.bcm)
		if addr != tt.address || op != tt.operation {
			t.Errorf("Test case %s failed. Expected: %v, %b. Got: %v, %b\n", tt.name, tt.address, tt.operation, addr, op)
		}
	}
}

func TestGetPinModePWM(t *testing.T) {
	c := new2837()

	testCases := []struct {
		name string
		in   int
		err  error
		out  PinMode
	}{
		{"p12", 12, nil, PinModeALT0},
		{"p18", 18, nil, PinModeALT5},
		{"p40", 40, nil, PinModeALT0},
		{"p52", 52, nil, PinModeALT1},
		{"p13", 13, nil, PinModeALT0},
		{"p19", 19, nil, PinModeALT5},
		{"p41", 41, nil, PinModeALT0},
		{"p45", 45, nil, PinModeALT0},
		{"p53", 53, nil, PinModeALT1},
		{"p10", 10, ErrNoPWM, PinModeNA},
	}

	for _, tt := range testCases {

		mode, err := c.getPinModePWM(tt.in)
		if mode != tt.out || err != tt.err {
			t.Errorf("Test case %s failed. Expected: %v, %v. Got: %v,%v\n", tt.name, tt.out, tt.err, mode, err)
		}
	}
}

func TestPWMClockCtl(t *testing.T) {
	testCases := []struct {
		name      string
		cfg       ClockConfig
		addr      uint64
		operation int
	}{
		{"clock_enab, mash 1, osc", ClockConfig{Enab: true, Mash: 1, Src: ClockSrcOsc}, 0x7E1010A0, 1509949969},
		{"clock_enab, mash 2, hdm1", ClockConfig{Enab: true, Mash: 2, Src: ClockSrcHDMI}, 0x7E1010A0, 1509950487},
	}

	c := new2837()
	for _, tt := range testCases {
		addr, _, operation := c.pwmClockCtl(tt.cfg)
		if addr != tt.addr || operation != tt.operation {
			t.Errorf("Test case %s failed. Expected: %v, %b. Got: %v,%b\n", tt.name, tt.addr, tt.operation, addr, operation)
		}
	}
}

func TestPWMClockDiv(t *testing.T) {
	testCases := []struct {
		name       string
		sourceFreq int
		freq       int
		addr       uint64
		operation  int
	}{
		{"65000 Hz", 19200000, 65000, 0x7E1010A4, 1511159335},
	}

	c := new2837()
	for _, tt := range testCases {
		addr, _, operation := c.pwmClockDiv(19200000, 65000)
		if addr != tt.addr || operation != tt.operation {
			t.Errorf("Test case %s failed. Expected: %v, %b. Got: %v,%b\n", tt.name, tt.addr, tt.operation, addr, operation)
		}
	}
}

func TestGpsel(t *testing.T) {
	c := new2837()

	testCases := []struct {
		name    string
		bcmNum  int
		pinMode PinMode

		address   uint64
		operation int
	}{
		{"p4 - in", 4, PinModeInput, 0x7E200000, 0},
		{"p9 - out", 9, PinModeOutput, 0x7E200000, 134217728},
	}

	for _, tt := range testCases {

		addr, _, op := c.gpgsel(tt.bcmNum, tt.pinMode)
		if addr != tt.address || op != tt.operation {
			t.Errorf("Test case %s failed. Expected: %v, %v. Got: %v,%v\n", tt.name, tt.address, tt.operation, addr, op)
		}
	}
}

func TestPwmCtl(t *testing.T) {

	c := new2837()

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

func TestPwmRng(t *testing.T) {
	c := new2837()

	address, _, op := c.pwmRng(12, 45)
	assert.Equal(t, address, uint64(0x7E20C010))
	assert.Equal(t, op, 45)

	address, _, op = c.pwmRng(120, 450)
	assert.Equal(t, address, uint64(0))
	assert.Equal(t, op, 450)

	address, _, op = c.pwmRng(19, 900)
	assert.Equal(t, address, uint64(0x7E20C020))
	assert.Equal(t, op, 900)
}

func TestPwmDat(t *testing.T) {
	c := new2837()

	address, _, op := c.pwmDat(12, 45)
	assert.Equal(t, address, uint64(0x7E20C014))
	assert.Equal(t, op, 45)

	address, _, op = c.pwmDat(120, 450)
	assert.Equal(t, address, uint64(0))
	assert.Equal(t, op, 450)

	address, _, op = c.pwmDat(19, 900)
	assert.Equal(t, address, uint64(0x7E20C024))
	assert.Equal(t, op, 900)
}
