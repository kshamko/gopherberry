package gopherberry

import (
	"errors"
	"fmt"
	//"github.com/kshamko/gopherberry/gpio"
	//"os"
)

type pinMode int

const (
	//NoBCMNum means that pin has no bcm number (ground, voltage pins)
	NoBCMNum = -1

	pinModeInput  pinMode = 0 //000
	pinModeOutput pinMode = 1 //001
	//pinModeALT0   pinMode = 4 //100
	//pinModeALT1   pinMode = 5 //101
	//pinModeALT2   pinMode = 6 //110
	//pinModeALT3   pinMode = 7 //111
	//pinModeALT4   pinMode = 3 //011
	//pinModeALT5   pinMode = 2 //010
)

var (
	//ErrNoPin error prodices when there is no bcm number for pin
	ErrNoPin = errors.New("no pin exists")
)

//Raspberry1 struct
type Raspberry struct {
	chip chip
	mmap *mmap
}

type chip interface {
	getPinBCM(pinNumBoard int) int
	getBaseVirtAddress() uint64
	getRegisters() map[string][]uint64

	gpgsel(bcm int, mode pinMode) (funcName string, addressOffset int, operation int)
	gpset(bcm int) (funcName string, addressOffset int, operation int)
}

//New func
func New() (*Raspberry, error) {
	chip := newChip2837()
	mmap, err := newMmap(chip.getRegisters(), chip.getBaseVirtAddress())
	if err != nil {
		fmt.Println("[ERROR] mmap err", err)
		return nil, err
	}

	return &Raspberry{
		chip: chip,
		mmap: mmap,
	}, nil
}

//GetPin returns pin object
func (r *Raspberry) GetPin(pinNumBoard int) (*Pin, error) {

	bcmNum := r.chip.getPinBCM(pinNumBoard)

	if bcmNum == NoBCMNum {
		return nil, ErrNoPin
	}

	return &Pin{
		bcmNum: bcmNum,
		chip:   r.chip,
		mmap:   r.mmap,
	}, nil

}
