package gopherberry

import "time"

func (p *Pin) SetFrequency(cfg ClockConfig, freq int) error {

	cfgToStop := cfg
	cfgToStop.Enab = false

	//stop clock
	addr, addrType, operation := p.pi.chip.clckCtl(p.bcmNum, cfgToStop)
	if addrType == addrBus {
		addr = p.pi.chip.addrBus2Phys(addr)
	}
	p.pi.mmapClock.run(addr, operation)
	//curCfg := p.pi.chip.clckCfg(p.pi.mmapClock.get(addr))

	//if curCfg.Busy { //wait until not busy
	time.Sleep(time.Microsecond * 10)

	//}
	//curClockCfg := p.pi.chip.clckCfg()

	addr1, addrType1, operation1 := p.pi.chip.clckDiv(p.bcmNum, freq)
	if addrType1 == addrBus {
		addr1 = p.pi.chip.addrBus2Phys(addr1)
	}

	p.pi.mmapClock.run(addr1, operation1)
	time.Sleep(time.Microsecond * 10) // ... so wait for them to take effect

	_, _, operation = p.pi.chip.clckCtl(p.bcmNum, cfg)
	return p.pi.mmapClock.run(addr, operation)
}

//
func (p *Pin) StopClock(cfg ClockConfig) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	cfg.Enab = false
	address, addressType, operation := p.pi.chip.clckCtl(p.bcmNum, cfg)
	return p.pi.runMmapClockCommand(address, addressType, operation)
}
