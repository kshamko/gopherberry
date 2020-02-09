package gopherberry

import (
	"sync"

	"github.com/pkg/errors"
)

type chipVersion int
type addressType int

const (
	//NoBCMNum means that pin has no bcm number (ground, voltage pins)
	NoBCMNum = -1
	//ARM2837 for corresonding chip type
	ARM2837 chipVersion = iota

	addrPhysical addressType = iota
	addrVirtual
	addrBus
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
//https://knowledge.ni.com/KnowledgeArticleDetails?id=kA00Z0000019OkFSAU - duty cycle descr
//https://electronics.stackexchange.com/questions/242293/is-there-an-ideal-pwm-frequency-for-dc-brush-motors - pwm freq
// https://www.precisionmicrodrives.com/content/ab-022-pwm-frequency-for-linear-motion-control/ - pwm freq
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
	mmapGPIO   *mmap
	mmapPWM    *mmap
	memOffsets map[uint64]int
	pwmRunning bool

	mu sync.Mutex
}

type chip interface {
	/*
		getBasePeriphialsAddressPhys() uint64
		//getBasePeriphialsAddressBus() uint64
		getGPIORegisters() gpioRegisters
		getPWMRegisters() pwmRegisters
		getPinModePWM(pinNumBCM int) (error, PinMode)*/

	getPinBCM(pinNumBoard int) int
	getGPIORegisters() (gpioRegisters, addressType)
	getPWMRegisters() (pwmRegisters, addressType)
	addrBus2Phys(uint64) uint64

	gpgsel(bcm int, mode PinMode) (registerAddress uint64, addressType addressType, operation int)
	gpset(bcm int) (registerAddress uint64, addressType addressType, operation int)
	gpclr(bcm int) (registerAddress uint64, addressType addressType, operation int)
	gplev(bcm int) (registerAddress uint64, addressType addressType, operation int)

	//pwmCtl(cfg1, cfg2 PWMChannelConfig) (registerAddress uint64, addressType addressType, operation int)
	//pwmRng() (registerAddress uint64, operation int)
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

	gpioMmap, err := raspberry.initMmapGPIO(c.getGPIORegisters())
	if err != nil {
		return nil, errors.Wrap(err, "can't init mmap")
	}
	raspberry.mmapGPIO = gpioMmap
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
/*func (r *Raspberry) StartPWM(cfg1, cfg2 PWMChannelConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	address, operation := r.chip.pwmCtl(cfg1, cfg2)
	offset, ok := r.memOffsets[address]
	if !ok {
		return ErrNoOffset
	}

	r.pwmRunning = true
	return r.mmap.run(offset, operation)
}*/

//StopPWM func disables PWM
func (r *Raspberry) StopPWM() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.pwmRunning = false

	return nil
}

/*func (r *Raspberry) initMmap(physicalAddresses []uint64) (*mmap, error) {
	sort.Slice(physicalAddresses, func(i, j int) bool { return physicalAddresses[i] < physicalAddresses[j] })
	return newMmap(physAddresses)
}*/

func (r *Raspberry) initMmapGPIO(gpioRegisters gpioRegisters, addressType addressType) (*mmap, error) {

	physicalAddresses := []uint64{}
	for _, registers := range gpioRegisters {
		for _, register := range registers {

			if addressType == addrBus {
				register = r.chip.addrBus2Phys(register)
			}
			physicalAddresses = append(physicalAddresses, register)
		}
	}
	return newMmap(physicalAddresses)
}

func (r *Raspberry) runMmapGPIOCommand(address uint64, addressType addressType, operation int) error {

	if addressType == addrBus {
		address = r.chip.addrBus2Phys(address)
	}

	return r.mmapGPIO.run(address, operation)
}

func (r *Raspberry) memStateGPIO(address uint64, addressType addressType) (int, error) {
	if addressType == addrBus {
		address = r.chip.addrBus2Phys(address)
	}
	return r.mmapGPIO.get(address)
}

/*func (r *Raspberry) initMmap() error {
	startPhysAddress := r.chip.getBasePeriphialsAddressPhys()
	mmap, err := newMmap(int64(startPhysAddress), os.Getpagesize())
	if err != nil {
		return err
	}
	r.mmap = mmap
	return nil
}

*/

//base := r.chip.getBasePeriphialsAddressPhys() & 0xff000000
/*physAddr := r.chip.addrBus2Phys(busAddress)
	offset := int(physAddr - r.chip.getBasePeriphialsAddressPhys())
	return r.mmap.run(offset, operation)
}*/

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
