package gopherberry

import (
	"errors"
	"fmt"
	"os/exec"
	"sync"
)

//Pin struct
type Pin struct {
	bcmNum int
	pi     *Raspberry
	mu     sync.Mutex

	curMode      pinMode
	edgeChan     chan EdgeType
	edgeToDetect EdgeType
	epoll        *Epoll
}

type EdgeType string

const (
	EdgeHigh EdgeType = "rising"
	EdgeLow  EdgeType = "falling"
	EdgeBoth EdgeType = "both"
	EdgeNone EdgeType = "none"
)

var (
	//ErrBadPinMode triggered when pin is not in a correct mode to call a function
	ErrBadPinMode = errors.New("pin is in the wrong mode")
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
	if p.curMode != pinModeOutput {
		return ErrBadPinMode
	}
	p.mu.Lock()
	defer p.mu.Unlock()

	address, operation := p.pi.chip.gpset(p.bcmNum)
	return p.runCommand(address, operation)
}

//SetLow sets an output to 0
func (p *Pin) SetLow() error {
	if p.curMode != pinModeOutput {
		return ErrBadPinMode
	}
	p.mu.Lock()
	defer p.mu.Unlock()

	address, operation := p.pi.chip.gpclr(p.bcmNum)
	return p.runCommand(address, operation)
}

//Level reports pin output state
func (p *Pin) Level() (bool, error) {
	if p.curMode != pinModeOutput {
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

//DetectEdge func
func (p *Pin) DetectEdge(edge EdgeType) (<-chan EdgeType, error) {
	if p.curMode != pinModeInput {
		return nil, ErrBadPinMode
	}
	p.mu.Lock()
	defer p.mu.Unlock()

	command := fmt.Sprintf("gpio edge %d %s", p.bcmNum, edge)
	_, err := exec.Command("/bin/bash", "-c", command).Output()
	if err != nil {
		return nil, err
	}

	if edge == EdgeNone {
		return nil, p.detectEdgeStop()
	}

	fileName := fmt.Sprintf("/sys/class/gpio/gpio%d/value", p.bcmNum)
	ep, err := NewEpoll(fileName)
	if err != nil {
		return nil, err
	}

	p.epoll = ep
	p.edgeChan = make(chan EdgeType)

	go func() {
		for {
			c := ep.Wait(SeekSet)
			data, ok := <-c
			if ok {

				if data[0] == 49 && (edge == EdgeBoth || edge == EdgeHigh) { //check 1
					p.edgeChan <- EdgeHigh
				}

				if data[0] == 48 && (edge == EdgeBoth || edge == EdgeLow) { //check 0
					p.edgeChan <- EdgeLow
				}
			} else {
				return
			}
		}
	}()

	return p.edgeChan, nil
}

//DetectEdgeStop stop
func (p *Pin) DetectEdgeStop() error {
	if p.curMode != pinModeInput {
		return ErrBadPinMode
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	_, err := p.DetectEdge(EdgeNone)

	return err
}

//
func (p *Pin) mode(mode pinMode) error {
	p.mu.Lock()
	defer p.mu.Unlock()

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

func (p *Pin) detectEdgeStop() (err error) {
	if p.edgeChan != nil {
		close(p.edgeChan)
	}
	if p.epoll != nil {
		err = p.epoll.Stop()
	}

	p.edgeToDetect = EdgeNone
	p.epoll = nil
	p.edgeChan = nil
	return err
}
