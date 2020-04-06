package gopherberry

import (
	"errors"
	"sync"
)

//Pin struct
type Pin struct {
	bcmNum int
	pi     *Raspberry
	mu     sync.Mutex

	curMode     PinMode
	isClockMode bool

	edgeChan     chan EdgeType
	edgeToDetect EdgeType
	epoll        *Epoll
}

type PinMode int
type EdgeType string

const (
	PinModeInput  PinMode = 0 //000
	PinModeOutput PinMode = 1 //001
	PinModeALT0   PinMode = 4 //100
	PinModeALT1   PinMode = 5 //101
	//PinModeALT2   PinMode = 6 //110
	//PinModeALT3   PinMode = 7 //111
	//PinModeALT4   PinMode = 3 //011
	PinModeALT5 PinMode = 2 //010

	PinModeNA PinMode = -1
)

var (
	//ErrBadPinMode triggered when pin is not in a correct mode to call a function
	ErrBadPinMode = errors.New("pin is in the wrong mode")
)

//GetMode func
func (p *Pin) GetMode() PinMode {
	return p.curMode
}

//
func (p *Pin) mode(mode PinMode, isClockMode bool) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.curMode = mode
	p.isClockMode = isClockMode

	address, addressType, operation := p.pi.chip.gpgsel(p.bcmNum, mode)
	return p.pi.memWriteGPIO(address, addressType, operation)
}
