package gopherberry

import (
	"errors"
	//"github.com/kshamko/gopherberry/gpio"
	//"os"
)

type pinMode int

const (
	//NoBCMNnm means that pin has no bcm number (ground, voltage pins)
	NoBCMNnm = -1

	pinModeInput  pinMode = 0 //000
	pinModeOutput pinMode = 1 //001
	pinModeALT0   pinMode = 4 //100
	pinModeALT1   pinMode = 5 //101
	pinModeALT2   pinMode = 6 //110
	pinModeALT3   pinMode = 7 //111
	pinModeALT4   pinMode = 3 //011
	pinModeALT5   pinMode = 2 //010
)

var (
	//ErrNoPin error prodices when there is no bcm number for pin
	ErrNoPin = errors.New("no pin exists")
)

//Chip struct
/*type Cho struct {
	//MemMap
	memMap *Mmap
	//PerBaseAddrPhys - periphial base physical address
	PerBaseAddrPhys uint64
	//PerBaseAddrBirt - periphial base virtual address
	PerBaseAddrVirt uint64
	//Board2BCM maps board pin number to "Broadcom SOC channel" number
	//https://raspberrypi.stackexchange.com/questions/12966/what-is-the-difference-between-board-and-bcm-for-gpio-pin-numbering
	Board2BCM map[int]int
	//GPIORegisters maps function to registers
	GPIORegisters map[string][]uint64
}*/

//Raspberry struct
type Raspberry struct {
	chip chip
}

type chip interface {
	getPinBCM(pinNumBoard int) int
	gpgsel(bcm int, mode pinMode) (addressOffset int, operation int)
}

//New func
func New() *Raspberry {
	return &Raspberry{
		chip: newChip2837(),
	}
}

//GetPin retuns pin object
func (r *Raspberry) GetPin(pinNumBoard int) (*Pin, error) {

	bcmNum := r.chip.getPinBCM(pinNumBoard)

	if bcmNum == NoBCMNnm {
		return nil, ErrNoPin
	}

	return &Pin{
		bcmNum: bcmNum,
		chip:   r.chip,
		//mMap:   c.memMap,
	}, nil

}
