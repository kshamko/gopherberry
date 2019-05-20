package gopherberry

import (
	"github.com/pkg/errors"
)

type pinMode int
type chipVersion int

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

	//ARM2837 for corresonding chip type
	ARM2837 chipVersion = iota

	addressInc = 4
)

var (
	//ErrNoPin error produced when there is no bcm number for pin
	ErrNoPin = errors.New("no pin exists")
	//ErrNoOffset error produced when there is mapping of address and offset in memOffsets
	ErrNoOffset = errors.New("no offset by address found")
	//ErrNoMmap produced when mmap data object has no offset
	ErrNoMmap = errors.New("No index in mmap")
)

//
type gpioRegisters map[string][]uint64

//Raspberry struct
type Raspberry struct {
	chip       chip
	mmap       *mmap
	memOffsets map[uint64]int
}

type chip interface {
	getPinBCM(pinNumBoard int) int
	getBasePeriphialsAddress() uint64
	getGPIORegisters() gpioRegisters

	gpgsel(bcm int, mode pinMode) (registerAddress uint64, operation int)
	gpset(bcm int) (registerAddress uint64, operation int)
	gpclr(bcm int) (registerAddress uint64, operation int)
	gplev(bcm int) (registerAddress uint64, operation int)
}

//New func
func New(chipVersion chipVersion) (*Raspberry, error) {

	c := newChip2837() //default
	if chipVersion == ARM2837 {
		c = newChip2837()
	}

	raspberry := &Raspberry{
		chip: c,
	}
	err := raspberry.initMmap()
	if err != nil {
		return nil, errors.Wrap(err, "can't init mmap")
	}
	return raspberry, nil
}

//GetPin returns pin object
func (r *Raspberry) GetPin(pinNumBoard int) (*Pin, error) {

	bcmNum := r.chip.getPinBCM(pinNumBoard)

	if bcmNum == NoBCMNum {
		return nil, ErrNoPin
	}

	return &Pin{
		bcmNum: bcmNum,
		pi:     r,
	}, nil
}

//
func (r *Raspberry) initMmap() error {

	var (
		minAddress, maxAddress uint64
	)

	for _, addresses := range r.chip.getGPIORegisters() {
		for _, address := range addresses {
			if minAddress == 0 {
				minAddress = address
			}

			if minAddress > address {
				minAddress = address
			}

			if maxAddress < address {
				maxAddress = address
			}
		}
	}

	mmapBaseAddr := mmapBaseAddress(minAddress, r.chip.getBasePeriphialsAddress())
	mmapLen := (maxAddress - minAddress) / addressInc

	mmap, err := newMmap(int64(mmapBaseAddr), int(mmapLen))

	if err != nil {
		return err
	}
	r.mmap = mmap
	r.memOffsets = offsets(r.chip.getGPIORegisters(), minAddress)

	return nil
}

//
func mmapBaseAddress(virtAddress, physBaseAddress uint64) uint64 {
	virtBase := virtAddress & 0xff000000
	return physBaseAddress + (virtAddress - virtBase)
}

//
func offsets(registers gpioRegisters, startAddress uint64) map[uint64]int {
	offsets := map[uint64]int{}
	for _, addresses := range registers {
		for _, address := range addresses {
			offsets[address] = offset(address, startAddress, addressInc)
		}
	}
	return offsets
}

//
func offset(address, startAddress uint64, addressInc int) int {
	return int((address - startAddress)) / addressInc
}
