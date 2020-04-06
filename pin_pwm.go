package gopherberry

//ModePWM set pin to PWM mode
func (p *Pin) ModePWM() error {

	mode, err := p.pi.chip.getPinModePWM(p.bcmNum)
	if err != nil {
		return err
	}
	return p.mode(mode, true)
}

//
func (p *Pin) DutyCycle(m, s int) error {
	a1, t1, op1 := p.pi.chip.pwmRng(p.bcmNum, s)
	a2, t2, op2 := p.pi.chip.pwmDat(p.bcmNum, m)

	p.pi.memWritePWM(a1, t1, op1)
	p.pi.memWritePWM(a2, t2, op2)

	return nil
}

//
func (p *Pin) SetMS() error {
	return nil
}
