package gopherberry

//ModePWM set pin to PWM mode
func (p *Pin) ModePWM() error {
	/*if !p.pi.pwmRunning {
		return ErrPWMStart
	}

	err, mode := p.pi.chip.getPinModePWM(p.bcmNum)
	if err != nil {
		return err
	}
	return p.mode(mode)*/
	return nil
}

func (p *Pin) SetMS() error {
	return nil
}
