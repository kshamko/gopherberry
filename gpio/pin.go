package gpio

type Pin struct {
}

func NewPin() (*Pin, error) {
	return nil, nil
}

//Select GPFSELn (?????)
func (p *Pin) Select() {

}

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
}


