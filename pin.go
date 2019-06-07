package gopherberry

import (
	"fmt"
	"os/exec"
)

//Pin struct
type Pin struct {
	bcmNum  int
	pi      *Raspberry
	curMode pinMode
}

type EdgeType string

const (
	EdgeHigh EdgeType = "rising"
	EdgeLow  EdgeType = "falling"
	EdgeBoth EdgeType = "both"
	EdgeNone EdgeType = "none"
)

//ModeInput sets pin to input mode
func (p *Pin) ModeInput() error {
	return p.mode(pinModeInput)
}

//ModeOutput sets pin to output mode
func (p *Pin) ModeOutput() error {
	return p.mode(pinModeOutput)
}

//GetMode func
//@todo implement
func (p *Pin) GetMode() pinMode {
	return p.curMode
}

//SetHigh sets an output to 1
func (p *Pin) SetHigh() error {
	address, operation := p.pi.chip.gpset(p.bcmNum)
	return p.runCommand(address, operation)
}

//SetLow sets an output to 0
func (p *Pin) SetLow() error {
	address, operation := p.pi.chip.gpclr(p.bcmNum)
	return p.runCommand(address, operation)
}

//Level reports pin output state
func (p *Pin) Level() (bool, error) {

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

//DetectEdge func
func (p *Pin) DetectEdge(edge EdgeType) (chan EdgeType, error) {
	command := fmt.Sprintf("gpio edge %d %s", p.bcmNum, edge)
	_, err := exec.Command("/bin/bash", "-c", command).Output()
	if err != nil {
		return nil, err
	}

	fileName := fmt.Sprintf("/sys/class/gpio/gpio%d/value", p.bcmNum)
	ep, err := NewEpoll(fileName)
	if err != nil {
		return nil, err
	}

	ch := make(chan EdgeType)

	go func() {
		for {
			c := ep.Wait()
			data, ok := <-c
			if ok {
				fmt.Println(data)

				if edge == EdgeBoth || edge == EdgeHigh { //check 1
					ch <- EdgeHigh
				}

				if edge == EdgeBoth || edge == EdgeLow { //check 0
					ch <- EdgeLow
				}
			} else {
				return
			}
		}
	}()

	return ch, nil
}

//
func (p *Pin) mode(mode pinMode) error {
	p.curMode = mode
	address, operation := p.pi.chip.gpgsel(p.bcmNum, mode)
	return p.runCommand(address, operation)
}

//
func (p *Pin) runCommand(address uint64, operation int) error {
	offset, ok := p.pi.memOffsets[address]
	if !ok {
		return ErrNoOffset
	}
	return p.pi.mmap.run(offset, operation)
}

//
func (p *Pin) memState(address uint64) (int, error) {
	offset, ok := p.pi.memOffsets[address]
	if !ok {
		return 0, ErrNoOffset
	}
	return p.pi.mmap.get(offset)
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
