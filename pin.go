package gopherberry

import "fmt"

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

func (p *Pin) mode(mode pinMode) error {

	addressOffset, operation := p.chip.gpgsel(p.bcmNum, mode)

	fmt.Println(addressOffset, operation)

	/*
		//calculate proper register offset
		registerOffset := p.BCMNum / 10 //1 register for 10 pins
		//calculate command. all commands are assumed to be 32-bit
		shift := (uint8(p.BCMNum) % 10) * 3 // 10 pins per register, command of 3 bits
		command := mode << shift
		p.curMode = mode
	*/

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
