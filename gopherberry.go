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
type registers map[string][]uint64

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
	chip      chip
	mmapGPIO  *mmap
	mmapPWM   *mmap
	mmapClock *mmap

	pwmRunning bool

	mu sync.Mutex
}

type chip interface {
	getPinBCM(pinNumBoard int) int
	getGPIORegisters() (registers, addressType)
	getPWMRegisters() (registers, addressType)
	getClockRegisters() (registers, addressType)

	addrBus2Phys(uint64) uint64
	getPinModePWM(pinNumBCM int) (PinMode, error)

	gpgsel(bcm int, mode PinMode) (registerAddress uint64, addressType addressType, operation int)
	gpset(bcm int) (registerAddress uint64, addressType addressType, operation int)
	gpclr(bcm int) (registerAddress uint64, addressType addressType, operation int)
	gplev(bcm int) (registerAddress uint64, addressType addressType, operation int)

	pwmCtl(cfg1, cfg2 PWMChannelConfig) (registerAddress uint64, addressType addressType, operation int)
	pwmRng(bcm int, val int) (registerAddress uint64, addressType addressType, operation int)
	pwmDat(bcm int, val int) (registerAddress uint64, addressType addressType, operation int)
	pwmClockCtl(cfg ClockConfig) (registerAddress uint64, addressType addressType, operation int)
	pwmClockDiv(sourceFreq, freq int) (registerAddress uint64, addressType addressType, operation int)

	//clckCtl(bcm int, cfg ClockConfig) (registerAddress uint64, addressType addressType, operation int)
	//clckDiv(bcm int, freq int) (registerAddress uint64, addressType addressType, operation int)
	//clckCfg(data int) ClockConfig
}

//New func
func New(chipVersion chipVersion) (*Raspberry, error) {

	c := new2837() //default
	if chipVersion == ARM2837 {
		c = new2837()
	}

	raspberry := &Raspberry{
		chip: c,
	}

	gpioMmap, err := raspberry.initMmap(c.getGPIORegisters())
	if err != nil {
		return nil, errors.Wrap(err, "can't init gpio mmap")
	}
	raspberry.mmapGPIO = gpioMmap

	pwmMmap, err := raspberry.initMmap(c.getPWMRegisters())
	if err != nil {
		return nil, errors.Wrap(err, "can't init pwm mmap")
	}
	raspberry.mmapPWM = pwmMmap

	clockMmap, err := raspberry.initMmap(c.getClockRegisters())
	if err != nil {
		return nil, errors.Wrap(err, "can't init clock mmap")
	}
	raspberry.mmapClock = clockMmap

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
	/*r.mu.Lock()
	defer r.mu.Unlock()
	address, addressType, operation := r.chip.pwmCtl(cfg1, cfg2)
	err := r.mmapPWM.run(r.physAddress(address, addressType), operation)
	if err == nil {
		r.pwmRunning = true
	}
	return err*/
	return nil
}

//StopPWM func disables PWM
func (r *Raspberry) StopPWM() error {
	/*r.mu.Lock()
	defer r.mu.Unlock()

	addr, addrType, operation := r.chip.pwmCtl(PWMChannelConfig{}, PWMChannelConfig{})
	return nil*/
	return nil

}

/**
* Memory mapped related functions
**/

func (r *Raspberry) initMmap(registers registers, addressType addressType) (*mmap, error) {

	physicalAddresses := []uint64{}
	for _, regs := range registers {
		for _, register := range regs {
			physicalAddresses = append(physicalAddresses, r.physAddress(register, addressType))
		}
	}
	return newMmap(physicalAddresses)
}

//
func (r *Raspberry) memWriteGPIO(address uint64, addressType addressType, operation int) error {
	return r.mmapGPIO.set(r.physAddress(address, addressType), operation)
}

//
func (r *Raspberry) memWritePWM(address uint64, addressType addressType, operation int) error {
	return r.mmapPWM.set(r.physAddress(address, addressType), operation)
}

//
func (r *Raspberry) memWriteClock(address uint64, addressType addressType, operation int) error {
	return r.mmapClock.set(r.physAddress(address, addressType), operation)
}

//
func (r *Raspberry) memReadGPIO(address uint64, addressType addressType) (int, error) {
	return r.mmapGPIO.get(r.physAddress(address, addressType))
}

//
func (r *Raspberry) memReadPWM(address uint64, addressType addressType) (int, error) {
	return r.mmapPWM.get(r.physAddress(address, addressType))
}

//
func (r *Raspberry) memReadClock(address uint64, addressType addressType) (int, error) {
	return r.mmapClock.get(r.physAddress(address, addressType))
}

func (r *Raspberry) physAddress(address uint64, addressType addressType) uint64 {
	if addressType == addrBus {
		return r.chip.addrBus2Phys(address)
	}

	return address
}
