package gopherberry

//ModePWM set pin to PWM mode
func (p *Pin) ModePWM() error {
	if !p.pi.pwmRunning {
		return ErrPWMStart
	}

	mode, err := p.pi.chip.getPinModePWM(p.bcmNum)
	if err != nil {
		return err
	}
	return p.mode(mode)
}

func (p *Pin) DutyCycle(m, s int) error {
	a1, t1, op1 := p.pi.chip.pwmRng(p.bcmNum, s)
	a2, t2, op2 := p.pi.chip.pwmDat(p.bcmNum, m)

	p.pi.runMmapPWMCommand(a1, t1, op1)
	p.pi.runMmapPWMCommand(a2, t2, op2)

	return nil
}

func (p *Pin) SetMS() error {
	return nil
}
