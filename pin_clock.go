package gopherberry

import (
	"time"

	"github.com/pkg/errors"
)

/*???????
* which reports the following:
* 0 0 Hz Ground
* 1 19.2 MHz oscillator
* 2 0 Hz testdebug0
* 3 0 Hz testdebug1
* 4 0 Hz PLLA
* 5 1000 MHz PLLC (changes with overclock settings)
* 6 500 MHz PLLD
* 7 216 MHz HDMI auxiliary
* 8-15 0 Hz Ground
 */

type ClockConfig struct {
	Mash int  //0,1,2,3
	Busy bool // read-only
	Enab bool //
	//0 = GND
	//1 = oscillator
	//2 = testdebug0/
	//3 = testdebug1
	//4 = PLLA per
	//5 = PLLC per
	//6 = PLLD per
	//7 = HDMI auxiliary
	Src ClockSource //
}

type ClockSource int

const (
	ClockSrcGND ClockSource = iota
	ClockSrcOsc
	ClockSrcTest1
	ClockSrcTest2
	ClockSrcPLLA
	ClockSrcPLLC
	ClockSrcPLLD
	ClockSrcHDMI
)

//@Todo:
//- verify busy flag for clock
//- check if pin in clock/pwm mode to make clock operations

//
func (p *Pin) SetFrequency(cfg ClockConfig, sourceFreq, freq int) error {

	if !p.isClockMode {
		return ErrBadPinMode
	}

	err := p.StopClock()
	if err != nil {
		return errors.Wrap(err, "SetFrequency")
	}

	//@todo check busy flag
	time.Sleep(time.Microsecond * 10)

	//If you want the cleanest clock source which is the XTAL (19.2MHz) crystal, then Clock source code = 0001b (oscilator)
	divRegister, divRegisterType, divOperation := p.pi.chip.pwmClockDiv(sourceFreq, freq)
	p.pi.memWriteClock(divRegister, divRegisterType, divOperation)

	time.Sleep(time.Microsecond * 10) // ... so wait for them to take effect

	err = p.StartClock(cfg)
	time.Sleep(time.Microsecond * 10)

	return errors.Wrap(err, "SetFrequency")
}

//
func (p *Pin) StartClock(cfg ClockConfig) error {
	if !p.isClockMode {
		return ErrBadPinMode
	}

	if !cfg.Enab {
		return errors.New("wrong config with Enab=false")
	}

	address, addressType, operation := p.pi.chip.pwmClockCtl(cfg)
	return p.pi.memWriteClock(address, addressType, operation)
}

//
func (p *Pin) StopClock() error {

	if !p.isClockMode {
		return ErrBadPinMode
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	cfg := ClockConfig{Enab: false}
	address, addressType, operation := p.pi.chip.pwmClockCtl(cfg)
	err := p.pi.memWriteClock(address, addressType, operation)

	if err != nil {
		return errors.Wrap(err, "StopClock")
	}
	//@todo check busy
	return nil
}
