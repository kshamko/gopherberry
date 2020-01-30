package gopherberry

import (
	"os"
	"sync"

	"github.com/pkg/errors"
)

type chipVersion int

const (
	//NoBCMNum means that pin has no bcm number (ground, voltage pins)
	NoBCMNum = -1
	//ARM2837 for corresonding chip type
	ARM2837    chipVersion = iota
	addressInc             = 1
)

var (
	//ErrNoPin error produced when there is no bcm number for the pin
	ErrNoPin = errors.New("no pin exists")
	//ErrNoOffset error produced when there is mapping of address and offset in memOffsets
	ErrNoOffset = errors.New("no offset by address found")
	//ErrNoMmap produced when mmap data object has no offset
	ErrNoMmap = errors.New("No index in mmap")
	//ErrNoPWM error produced when PWM is not available for the pin
	ErrNoPWM = errors.New("PWM mode is not available for the pin")
	//ErrPWMStart error produced if PWM was not started by StartPWM()
	ErrPWMStart = errors.New("PWM not started")
)

//
type gpioRegisters map[string][]uint64
type pwmRegisters map[string]uint64

//PWMChannelConfig struct represents configuration of a PWM channel
//Pi has 2 PWM channels
//https://github.com/RichardChambers/raspberrypi/wiki/Notes-on:-Pulse-Width-Modulation-(PWM)---Discussion
//
type PWMChannelConfig struct {
	//0 - balanced mode (pulse-density)
	//1 - mark-space mode (M/S)
	MSEnable   int
	UseFIF0    int
	Polarity   int
	SilenceBit int
	RepeatLast int
	//https://www.raspberrypi.org/forums/viewtopic.php?t=16181
	// if set to 1 than MSEnable has no effect
	// 0 - PWM mode
	// 1 - serializer mode. Data is transmitted MSB first and truncated or zero-padded depending on PWM_RNGi.
	Mode        int
	ChanEnabled int
}

//Raspberry struct
type Raspberry struct {
	chip       chip
	mmap       *mmap
	memOffsets map[uint64]int
	pwmRunning bool

	mu sync.Mutex
}

type chip interface {
	getPinBCM(pinNumBoard int) int
	getBasePeriphialsAddressPhys() uint64
	getBasePeriphialsAddressBus() uint64
	getGPIORegisters() gpioRegisters
	getPWMRegisters() pwmRegisters
	getPinModePWM(pinNumBCM int) (error, PinMode)

	gpgsel(bcm int, mode PinMode) (registerAddress uint64, operation int)
	gpset(bcm int) (registerAddress uint64, operation int)
	gpclr(bcm int) (registerAddress uint64, operation int)
	gplev(bcm int) (registerAddress uint64, operation int)

<<<<<<< HEAD
	pwmCtl(cfg1, cfg2 PWMChannelConfig) (registerAddress uint64, operation int)
=======
	pwmAltFunc(bcm int) (alt PinMode, pwmChanNum int, err error)
>>>>>>> 73d75e7cced37c7de088a06c15bda7037fe81722
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

//StartPWM func enables PWM
//
//
func (r *Raspberry) StartPWM(cfg1, cfg2 PWMChannelConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	address, operation := r.chip.pwmCtl(cfg1, cfg2)
	offset, ok := r.memOffsets[address]
	if !ok {
		return ErrNoOffset
	}

	r.pwmRunning = true
	return r.mmap.run(offset, operation)
}

//StopPWM func disables PWM
func (r *Raspberry) StopPWM() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.pwmRunning = false

	return nil
}

func (r *Raspberry) initMmap() error {
	startPhysAddress := r.chip.getBasePeriphialsAddressPhys()
	mmap, err := newMmap(int64(startPhysAddress), os.Getpagesize())
	if err != nil {
		return err
	}
	r.mmap = mmap
	return nil
}

func (r *Raspberry) runMmapCommand(busAddress uint64, operation int) error {
	/*offset, ok := p.pi.memOffsets[address]
	if !ok {
		return ErrNoOffset
	}*/
	//base := r.chip.getBasePeriphialsAddressPhys() & 0xff000000

	offset := int(busAddress - r.chip.getBasePeriphialsAddressBus())
	return r.mmap.run(offset, operation)
}

/*func (r *Raspberry) Close() {
	r.mmap.Close()
}*/

//TODO refactor mmap load
/*func (r *Raspberry) initMmap() error {

/*var (
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
}*/

/*minAddress := uint64(0x7E200000)
	mmapBaseAddr := mmapBaseAddress(minAddress, r.chip.getBasePeriphialsAddress())
	mmapLen := os.Getpagesize() //(maxAddress - minAddress) / addressInc

	mmap, err := newMmap(int64(mmapBaseAddr), mmapLen)

	if err != nil {
		return err
	}
	r.mmap = mmap
	r.memOffsets = offsetsGPIO(r.chip.getGPIORegisters(), minAddress)
	for k, v := range offsetsPWM(r.chip.getPWMRegisters(), minAddress) {
		r.memOffsets[k] = v
	}

	return nil
}

//
func mmapBaseAddress(virtAddress, physBaseAddress uint64) uint64 {
	virtBase := virtAddress & 0xff000000
	return physBaseAddress + (virtAddress - virtBase)
}

//
func offsetsGPIO(registers gpioRegisters, startAddress uint64) map[uint64]int {
	offsets := map[uint64]int{}
	for _, addresses := range registers {
		for _, address := range addresses {
			offsets[address] = offset(address, startAddress, addressInc)
		}
	}
	return offsets
}

func offsetsPWM(registers pwmRegisters, startAddress uint64) map[uint64]int {
	offsets := map[uint64]int{}
	for _, address := range registers {
		offsets[address] = offset(address, startAddress, addressInc)
	}
	return offsets
}

//
func offset(address, startAddress uint64, addressInc int) int {
	return int((address - startAddress)) / addressInc
}*/
