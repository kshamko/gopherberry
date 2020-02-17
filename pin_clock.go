package gopherberry

import (
	"time"
)

//@Todo:
//- verify busu flag for clock
//- check if pin in clock/pwm mode to make clock operations

//
func (p *Pin) SetFrequency(cfg ClockConfig, freq int) error {

	/*err := p.StopClock()
	if err != nil {
		return nil
	}

	//@todo check busy flag
	time.Sleep(time.Microsecond * 10)*/

	cfg.Enab = false
	p.StartClock(cfg)

	addr1, addrType1, operation1 := p.pi.chip.clckDiv(p.bcmNum, freq)
	if addrType1 == addrBus {
		addr1 = p.pi.chip.addrBus2Phys(addr1)
	}
	//
	//divi = 300, divf = 0 freq=64000
	p.pi.mmapClock.run(addr1, operation1)
	time.Sleep(time.Microsecond * 10) // ... so wait for them to take effect

	cfg.Enab = true
	err := p.StartClock(cfg)
	time.Sleep(time.Microsecond * 10)
	/*
		_, _, operation = p.pi.chip.clckCtl(p.bcmNum, cfg)
		return p.pi.mmapClock.run(addr, operation)*/

	return err
}

//
func (p *Pin) StartClock(cfg ClockConfig) error {
	if !cfg.Enab {
		//return errors.New("wrong config with Enab=false")
	}
	address, addressType, operation := p.pi.chip.clckCtl(p.bcmNum, cfg)
	return p.pi.runMmapClockCommand(address, addressType, operation)
}

//
func (p *Pin) StopClock() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	cfg := ClockConfig{Enab: false}
	address, addressType, operation := p.pi.chip.clckCtl(p.bcmNum, cfg)
	return p.pi.runMmapClockCommand(address, addressType, operation)
}
