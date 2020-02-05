package gopherberry

//ModeOutput sets pin to output mode
func (p *Pin) ModeOutput() error {
	return p.mode(PinModeOutput)
}

//SetHigh sets an output to 1
func (p *Pin) SetHigh() error {
	if p.curMode != PinModeOutput {
		return ErrBadPinMode
	}
	p.mu.Lock()
	defer p.mu.Unlock()

	address, operation := p.pi.chip.gpset(p.bcmNum)
	return p.pi.runMmapCommand(address, operation)
}

//SetLow sets an output to 0
func (p *Pin) SetLow() error {
	if p.curMode != PinModeOutput {
		return ErrBadPinMode
	}
	p.mu.Lock()
	defer p.mu.Unlock()

	address, operation := p.pi.chip.gpclr(p.bcmNum)
	return p.pi.runMmapCommand(address, operation)
}

//Level reports pin output state
func (p *Pin) Level() (bool, error) {
	if p.curMode != PinModeOutput {
		return false, ErrBadPinMode
	}
	p.mu.Lock()
	defer p.mu.Unlock()

	address, operation := p.pi.chip.gplev(p.bcmNum)
	state, err := p.memState(address)
	if err != nil {
		return false, err
	}

	if state&operation == 0 {
		return false, nil
	}

	return true, nil
}
