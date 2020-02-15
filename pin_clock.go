package gopherberry

import "time"

func (p *Pin) SetFrequency(freq int) error {
	// TODO: would be nice to choose best clock source depending on target frequency, oscilator is used for now
	const sourceFreq = 19200000 // oscilator frequency
	const divMask = 4095        // divi and divf have 12 bits each

	divi := uint32(sourceFreq / freq)
	divf := uint32(((sourceFreq % freq) << 12) / freq)

	divi &= divMask
	divf &= divMask

	/*clkCtlReg := 28
	  clkDivReg := 28
	  switch pin {
	  case 4, 20, 32, 34: // clk0
	  	clkCtlReg += 0
	  	clkDivReg += 1
	  case 5, 21, 42, 44: // clk1
	  	clkCtlReg += 2
	  	clkDivReg += 3
	  case 6, 43: // clk2
	  	clkCtlReg += 4
	  	clkDivReg += 5
	  case 12, 13, 40, 41, 45, 18, 19: // pwm_clk - shared clk for both pwm channels
	  	clkCtlReg += 12
	  	clkDivReg += 13
	  	StopPwm() // pwm clk busy wont go down without stopping pwm first
	  	defer StartPwm()
	  default:
	  	return
	  }*/

	mash := uint32(1 << 9) // 1-stage MASH
	if divi < 2 || divf == 0 {
		mash = 0
	}

	const PASSWORD = 0x5A000000
	const busy = 1 << 7
	const enab = 1 << 4
	const src = 1 << 0 // oscilator

	clkMem[clkCtlReg] = PASSWORD | (clkMem[clkCtlReg] &^ enab) // stop gpio clock (without changing src or mash)
	for clkMem[clkCtlReg]&busy != 0 {
		time.Sleep(time.Microsecond * 10)
	} // ... and wait for not busy

	clkMem[clkCtlReg] = PASSWORD | mash | src          // set mash and source (without enabling clock)
	clkMem[clkDivReg] = PASSWORD | (divi << 12) | divf // set dividers

	// mash and src can not be changed in same step as enab, to prevent lock-up and glitches
	time.Sleep(time.Microsecond * 10) // ... so wait for them to take effect

	clkMem[clkCtlReg] = PASSWORD | mash | src | enab // finally start clock

	// NOTE without root permission this changes will simply do nothing successfully
}

//
func (p *Pin) StopClock() error {
	address, addressType, operation := p.pi.chip.clckCtl(p.bcmNum, false)
	return p.pi.runMmapClockCommand(address, addressType, operation)
}
