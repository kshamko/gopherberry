package gopherberry

//Pin struct
type Pin struct {
	bcmNum  int
	chip    chip
	mmap    *mmap
	curMode pinMode
}

//ModeInput sets pin to input mode
func (p *Pin) ModeInput() error {
	return p.mode(pinModeInput)
}

//ModeOutput sets pin to output mode
func (p *Pin) ModeOutput() error {
	return p.mode(pinModeOutput)
}

//SetHigh sets an output to 1
func (p *Pin) SetHigh() error {
	funcName, addressOffset, operation := p.chip.gpset(p.bcmNum)
	return p.mmap.run(funcName, addressOffset, operation)
}

//
func (p *Pin) mode(mode pinMode) error {
	p.curMode = mode
	funcName, addressOffset, operation := p.chip.gpgsel(p.bcmNum, mode)
	return p.mmap.run(funcName, addressOffset, operation)
}

/*
//SetOutput GPSETn (R/W)
func (p *Pin) SetOutput() error {
	return nil
}

//ClearOutput GPCLRn (R/W)
func (p *Pin) ClearOutput() error {
	return nil
}

//Level GPLEVn (R/W)
func (p *Pin) Level() (bool, error) {
	return false, nil
}

//DetectStatusEvent GPEDSn (R/W)
func (p *Pin) DetectStatusEvent() error {
	return nil
}

//DetectEdgeRising (GPRENn) (R/W)
func (p *Pin) DetectEdgeRising() error {
	return nil
}

//DetectEdgeFalling (GPRENn) (R/W)
func (p *Pin) DetectEdgeFalling() error {
	return nil
}

//HighDetectEnable (GPHENn)
func (p *Pin) HighDetectEnable() error {
	return nil
}

//LowDetectEnable (GPLENn)
func (p *Pin) LowDetectEnable() error {
	return nil
}

//DetectEdgeRisingAsync (GPARENn)
func (p *Pin) DetectEdgeRisingAsync() error {
	return nil
}

//DetectEdgFallingAsync (GPAFENn)
func (p *Pin) DetectEdgFallingAsync() error {
	return nil
}*/
