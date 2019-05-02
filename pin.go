package gopherberry

type pinMode int

//See GPFSELn spec for details
const (
	pinModeInput  pinMode = 0 //000
	pinModeOutput pinMode = 1 //001

	pinModeALT0 pinMode = 4 //100
	pinModeALT1 pinMode = 5 //101
	pinModeALT2 pinMode = 6 //110
	pinModeALT3 pinMode = 7 //111
	pinModeALT4 pinMode = 3 //011
	pinModeALT5 pinMode = 2 //010
)

//Pin struct
type Pin struct {
	BCMNum  int
	mMap    *Mmap
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

//See GPFSEL(0..5) spec in ARM datasheet for details. GPSEL uses 6 registers
//GPSET0 for pins 0-9
//GPSET1 for pins 10-19
//GPSET2 for pins 20-29
//GPSET3 for pins 30-39
//GPSET4 for pins 40-49
//GPSET5 for pins 50-53
func (p *Pin) mode(mode pinMode) error {

	//calculate proper register offset
	registerOffset := p.BCMNum / 10 //1 register for 10 pins
	//calculate command. all commands are assumed to be 32-bit
	shift := (uint8(p.BCMNum) % 10) * 3 // 10 pins per register, command of 3 bits
	command := mode << shift
	p.curMode = mode

	return nil
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
