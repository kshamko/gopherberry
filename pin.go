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

	curMode      PinMode
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
//@todo implement
func (p *Pin) GetMode() PinMode {
	return p.curMode
}

//
func (p *Pin) mode(mode PinMode) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.curMode = mode
	address, operation := p.pi.chip.gpgsel(p.bcmNum, mode)
	return p.pi.runMmapCommand(address, operation)
}

//
/*func (p *Pin) runCommand(address uint64, operation int) error {
	offset, ok := p.pi.memOffsets[address]
	if !ok {
		return ErrNoOffset
	}
	return p.pi.mmap.run(offset, operation)
}*/

//
func (p *Pin) memState(address uint64) (int, error) {
	offset, ok := p.pi.memOffsets[address]
	if !ok {
		return 0, ErrNoOffset
	}
	return p.pi.mmap.get(offset)
}
