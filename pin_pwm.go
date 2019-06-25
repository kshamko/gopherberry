package gopherberry

//ModePWM set pin to PWM mode
func (p *Pin) ModePWM() error {
	m, _, err := p.pi.chip.pwmAltFunc(p.bcmNum)
	if err != nil {
		return err
	}
	return p.mode(m)
}
