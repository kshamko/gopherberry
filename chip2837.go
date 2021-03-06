package gopherberry

//Chip2837 implementation for raspberry 3+
type Chip2837 struct {
	periphialsBaseAddrPhys uint64
	//periphialsBaseAddrVirt uint64
	periphialsBaseAddrBus uint64

	//Board2BCM maps board pin number to "Broadcom SOC channel" number
	//nolint https://raspberrypi.stackexchange.com/questions/12966/what-is-the-difference-between-board-and-bcm-for-gpio-pin-numbering
	board2BCM map[int]int
	//GPIORegisters maps function to registers
	gpioRegisters  registers
	clockRegisters registers
	pwmRegisters   registers

	pwm0, pwm1             map[int]PinMode
	clock0, clock1, clock2 map[int]PinMode
}

const clockPassword = 0x5A000000

//New() func
func new2837() chip {

	//Peripherals (at physical address 0x3F000000 on) are mapped into the kernel virtual address
	//space starting at address 0xF2000000. Thus a peripheral advertised here at bus address
	//0x7Ennnnnn is available in the ARM kenel at virtual address 0xF2nnnnnn.
	chip := &Chip2837{
		// I/O devices are assigned physical memory addresses, which the Linux kernel prevents user programs from accessing.
		periphialsBaseAddrPhys: 0x3F000000,
		// We need to map this physical memory into the program's virtual addressing scheme so that the program can access it.
		//periphialsBaseAddrVirt: 0xF2000000,
		periphialsBaseAddrBus: 0x7E000000,
		board2BCM: map[int]int{
			1:  NoBCMNum, //3v3 power
			2:  NoBCMNum, //5v power
			3:  2,        //SDA
			4:  NoBCMNum, //5v power
			5:  3,        //SCL
			6:  NoBCMNum, //ground
			7:  4,        //GPCLK0
			8:  14,       //TXD
			9:  NoBCMNum, //ground
			10: 15,       //RXD
			11: 17,
			12: 18, //PWM0
			13: 27,
			14: NoBCMNum, //ground
			15: 22,
			16: 23,
			17: NoBCMNum, //3v3 power
			18: 24,
			19: 10,       //MOSI
			20: NoBCMNum, //ground
			21: 9,        //MISO
			22: 25,
			23: 11,       //SCLK
			24: 8,        //CE0
			25: NoBCMNum, //ground
			26: 7,        //CE1
			27: 0,        //ID_SD
			28: 1,        //ID_SC
			29: 5,
			30: NoBCMNum, //ground
			31: 6,
			32: 12, //PWM0
			33: 13, //PWM1
			34: 16,
			35: 19, //MISO
			36: 16,
			37: 26,
			38: 20,       //MOSI
			39: NoBCMNum, //ground
			40: 21,       //SCLK
		},
		gpioRegisters: map[string][]uint64{
			"GPFSEL":   {0x7E200000, 0x7E200004, 0x7E200008, 0x7E20000C, 0x7E200010, 0x7E200014}, //select pin mode
			"rsrvd1":   {0x7E200018},
			"GPSET":    {0x7E20001C, 0x7E200020}, //set pin high
			"rsrvd2":   {0x7E200024},
			"GPCLR":    {0x7E200028, 0x7E20002C}, //set pin low
			"rsrvd3":   {0x7E200030},
			"GPLEV":    {0x7E200034, 0x7E200038}, //get pin level
			"rsrvd4":   {0x7E20003C},
			"GPEDS":    {0x7E200040, 0x7E200044}, //rw
			"rsrvd5":   {0x7E200048},
			"GPREN":    {0x7E20004C, 0x7E200050}, //rw
			"rsrvd6":   {0x7E200054},
			"GPFEN":    {0x7E200058, 0x7E20005C}, //rw
			"rsrvd7":   {0x7E200060},
			"GPHEN":    {0x7E200064, 0x7E200068}, //rw
			"rsrvd8":   {0x7E20006C},
			"GPLEN":    {0x7E200070, 0x7E200074}, //rw
			"rsrvd9":   {0x7E200078},
			"GPAREN":   {0x7E20007C, 0x7E200080}, //rw
			"rsrvd10":  {0x7E200084},
			"GPAFEN":   {0x7E200088, 0x7E20008C}, //rw
			"rsrvd11":  {0x7E200090},
			"GPPUD":    {0x7E200094},             //rw
			"GPPUDCLK": {0x7E200098, 0x7E20009C}, //rw
		},
		pwm0: map[int]PinMode{ //map[bcmNum]PinMode
			12: PinModeALT0,
			18: PinModeALT5,
			40: PinModeALT0,
			52: PinModeALT1,
		},
		pwm1: map[int]PinMode{
			13: PinModeALT0,
			19: PinModeALT5,
			41: PinModeALT0,
			45: PinModeALT0,
			53: PinModeALT1,
		},
		pwmRegisters: map[string][]uint64{
			"CTL":   {0x7E20C000},
			"STA":   {0x7E20C004},
			"DMAC":  {0x7E20C008},
			"RSRV0": {0x7E20C00C},
			"RNG1":  {0x7E20C010},
			"DAT1":  {0x7E20C014},
			"FIF1":  {0x7E20C018},
			"RSRV1": {0x7E20C01C},
			"RNG2":  {0x7E20C020},
			"DAT2":  {0x7E20C024},
		},
		clock0: map[int]PinMode{},
		clock1: map[int]PinMode{},
		clock2: map[int]PinMode{},
		//https://www.scribd.com/doc/127599939/BCM2835-Audio-clocks#scribd
		clockRegisters: map[string][]uint64{
			"Unknown":  {0x7e101000, 0x7e101004, 0x7e101008, 0x7e10100c, 0x7e101010, 0x7e101014, 0x7e101018, 0x7e10101c, 0x7e101020, 0x7e101024, 0x7e101028, 0x7e10102c, 0x7e101030, 0x7e101034, 0x7e101038, 0x7e10103c, 0x7e101040, 0x7e101044, 0x7e101048, 0x7e10104c, 0x7e101050, 0x7e101054, 0x7e101058, 0x7e10105c, 0x7e101060, 0x7e101064, 0x7e101068, 0x7e10106c},
			"GP0CTL":   {0x7E101070},
			"GP0DIV":   {0x7E101074},
			"GP1CTL":   {0x7E101078},
			"GP1DIV":   {0x7E10107C},
			"GP2CTL":   {0x7E101080},
			"GP2DIV":   {0x7E101084},
			"Unknown1": {0x7e101088, 0x7e10108c, 0x7e101090, 0x7e101094},
			"PCMCTL":   {0x7E101098},
			"PCMDIV":   {0x7E10109C},
			"PWMCTL":   {0x7E1010A0},
			"PWMDIV":   {0x7E1010A4},
		},
	}

	return chip
}

func (chip *Chip2837) addrBus2Phys(addr uint64) uint64 {
	return addr - chip.periphialsBaseAddrBus + chip.periphialsBaseAddrPhys
}

//
func (chip *Chip2837) getGPIORegisters() (registers, addressType) {
	return chip.gpioRegisters, addrBus
}

func (chip *Chip2837) getPWMRegisters() (registers, addressType) {
	return chip.pwmRegisters, addrBus
}

func (chip *Chip2837) getClockRegisters() (registers, addressType) {
	return chip.clockRegisters, addrBus
}

func (chip *Chip2837) getPinBCM(pinNumBoard int) int {
	bcm, ok := chip.board2BCM[pinNumBoard]
	if ok {
		return bcm
	}

	return NoBCMNum
}

func (chip *Chip2837) getPinModePWM(pinNumBCM int) (PinMode, error) {

	if val, ok := chip.pwm0[pinNumBCM]; ok {
		return val, nil
	}
	if val, ok := chip.pwm1[pinNumBCM]; ok {
		return val, nil
	}
	return PinModeNA, ErrNoPWM
}

//
func (chip *Chip2837) gpgsel(bcm int, mode PinMode) (registerAddress uint64, addressType addressType, operation int) {
	//calculate proper register offset
	addressOffset := bcm / 10 //1 register for 10 pins
	//calculate operation. all operations are assumed to be 32-bit
	shift := (uint8(bcm) % 10) * 3 // 10 pins per register, command of 3 bits
	operation = int(mode) << shift
	return chip.gpioRegisters["GPFSEL"][addressOffset], addrBus, operation
}

//
func (chip *Chip2837) gpset(bcm int) (registerAddress uint64, addressType addressType, operation int) {
	return chip.twoBankCommand(bcm, "GPSET", addrBus)
}

//
func (chip *Chip2837) gpclr(bcm int) (registerAddress uint64, addressType addressType, operation int) {
	return chip.twoBankCommand(bcm, "GPCLR", addrBus)
}

//
func (chip *Chip2837) gplev(bcm int) (registerAddress uint64, addressType addressType, operation int) {
	return chip.twoBankCommand(bcm, "GPLEV", addrBus)
}

func (chip *Chip2837) pwmCtl(cfg1, cfg2 PWMChannelConfig) (registerAddress uint64, addressType addressType, operation int) {
	operation = cfg2.MSEnable<<15 + cfg2.UseFIF0<<13 + cfg2.Polarity<<12 + cfg2.SilenceBit<<11
	operation += cfg2.RepeatLast<<10 + cfg2.Mode<<9 + cfg2.ChanEnabled<<8
	operation += cfg1.MSEnable<<7 + cfg1.UseFIF0<<5 + cfg1.Polarity<<4 + cfg1.SilenceBit<<3
	operation += cfg1.RepeatLast<<2 + cfg1.Mode<<1 + cfg1.ChanEnabled
	return chip.pwmRegisters["CTL"][0], addrBus, operation
}

func (chip *Chip2837) pwmRng(bcm int, val int) (registerAddress uint64, addressType addressType, operation int) {
	if _, ok := chip.pwm0[bcm]; ok {
		return chip.pwmRegisters["RNG1"][0], addrBus, val
	}

	if _, ok := chip.pwm1[bcm]; ok {
		return chip.pwmRegisters["RNG2"][0], addrBus, val
	}

	return 0, addrBus, val
}

func (chip *Chip2837) pwmDat(bcm int, val int) (registerAddress uint64, addressType addressType, operation int) {
	if _, ok := chip.pwm0[bcm]; ok {
		return chip.pwmRegisters["DAT1"][0], addrBus, val
	}

	if _, ok := chip.pwm1[bcm]; ok {
		return chip.pwmRegisters["DAT2"][0], addrBus, val
	}

	return 0, addrBus, val
}

//
func (chip *Chip2837) pwmClockCtl(cfg ClockConfig) (registerAddress uint64, addressType addressType, operation int) {
	return chip.clckCtl("PWMCTL", cfg)
}

//
func (chip *Chip2837) pwmClockDiv(sourceFreq, freq int) (registerAddress uint64, addressType addressType, operation int) {
	return chip.clckDiv("PWMDIV", sourceFreq, freq)
}

/*
HELPERS
*/

//
func (chip *Chip2837) clckCtl(registerName string, cfg ClockConfig) (registerAddress uint64, addressType addressType, operation int) {
	enab := 1 << 4
	disab := 0 << 4
	operation = clockPassword | cfg.Mash<<9 | int(cfg.Src) | disab
	if cfg.Enab {
		operation = operation | enab
	}
	return chip.clockRegisters[registerName][0], addrBus, operation
}

//If you want the cleanest clock source which is the XTAL (19.2MHz) crystal, then Clock source code = 0001b (oscilator)
func (chip *Chip2837) clckDiv(registerName string, sourceFreq, freq int) (registerAddress uint64, addressType addressType, operation int) {
	//freq = source / ( DIVI + DIVF / 4096 ) - from the reference doc
	divi := sourceFreq / freq
	divr := sourceFreq % freq
	divf := int(float32(divr) * 4096.0 / float32(freq))
	return chip.clockRegisters[registerName][0], addrBus, (clockPassword | (divi << 12) | divf)
}

//
func (chip *Chip2837) twoBankCommand(bcm int, commandName string, addrType addressType) (registerAddress uint64, addressType addressType, operation int) {
	addressOffset := bcm / 32 //1 register for 32 pins
	shift := (uint8(bcm) % 32)
	operation = 1 << shift
	return chip.gpioRegisters[commandName][addressOffset], addrType, operation
}
