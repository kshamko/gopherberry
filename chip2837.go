package gopherberry

//Chip2837 implementation for raspberry 3+
type Chip2837 struct {
	periphialsBaseAddrPhys uint64
	//periphialsBaseAddrVirt uint64
	//periphialsBaseAddrBus  uint64

	//Board2BCM maps board pin number to "Broadcom SOC channel" number
	//nolint https://raspberrypi.stackexchange.com/questions/12966/what-is-the-difference-between-board-and-bcm-for-gpio-pin-numbering
	board2BCM map[int]int
	//GPIORegisters maps function to registers
	gpioRegisters gpioRegisters
}

//NewChip2837 func
func newChip2837() chip {

	//Peripherals (at physical address 0x3F000000 on) are mapped into the kernel virtual address
	//space starting at address 0xF2000000. Thus a peripheral advertised here at bus address
	//0x7Ennnnnn is available in the ARM kenel at virtual address 0xF2nnnnnn.
	chip := &Chip2837{
		// I/O devices are assigned physical memory addresses, which the Linux kernel prevents user programs from accessing.
		periphialsBaseAddrPhys: 0x3F000000,
		// We need to map this physical memory into the program's virtual addressing scheme so that the program can access it.
		//periphialsBaseAddrVirt: 0xF2000000,
		//periphialsBaseAddrBus:  0x7E000000,
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
			"GPFSEL":   {0x7E200000, 0x7E200004, 0x7E200008, 0x7E20000C, 0x7E200010, 0x7E200014}, //rw
			"GPSET":    {0x7E20001C, 0x7E200020},                                                 // w
			"GPCLR":    {0x7E200028, 0x7E20002C},                                                 // w
			"GPLEV":    {0x7E200034, 0x7E200038},                                                 // r
			"GPEDS":    {0x7E200040, 0x7E200044},                                                 // rw
			"GPREN":    {0x7E20004C, 0x7E200050},                                                 //rw
			"GPFEN":    {0x7E200058, 0x7E20005C},                                                 //rw
			"GPHEN":    {0x7E200064, 0x7E200068},                                                 //rw
			"GPLEN":    {0x7E200070, 0x7E200074},                                                 //rw
			"GPAREN":   {0x7E20007C, 0x7E200080},                                                 //rw
			"GPAFEN":   {0x7E200088, 0x7E20008C},                                                 //rw
			"GPPUD":    {0x7E200094},                                                             //rw
			"GPPUDCLK": {0x7E200098, 0x7E20009C},                                                 //rw
		},
	}

	return chip
}

//
func (chip *Chip2837) getBasePeriphialsAddress() uint64 {
	return chip.periphialsBaseAddrPhys
}

//
func (chip *Chip2837) getGPIORegisters() gpioRegisters {
	return chip.gpioRegisters
}

func (chip *Chip2837) getPinBCM(pinNumBoard int) int {
	bcm, ok := chip.board2BCM[pinNumBoard]
	if ok {
		return bcm
	}

	return NoBCMNum
}

//
func (chip *Chip2837) gpgsel(bcm int, mode PinMode) (registerAddress uint64, operation int) {
	//calculate proper register offset
	addressOffset := bcm / 10 //1 register for 10 pins
	//calculate operation. all operations are assumed to be 32-bit
	shift := (uint8(bcm) % 10) * 3 // 10 pins per register, command of 3 bits
	operation = int(mode) << shift
	return chip.gpioRegisters["GPFSEL"][addressOffset], operation
}

//
func (chip *Chip2837) gpset(bcm int) (registerAddress uint64, operation int) {
	return chip.twoBankCommand(bcm, "GPSET")
}

//
func (chip *Chip2837) gpclr(bcm int) (registerAddress uint64, operation int) {
	return chip.twoBankCommand(bcm, "GPCLR")
}

//
func (chip *Chip2837) gplev(bcm int) (registerAddress uint64, operation int) {
	return chip.twoBankCommand(bcm, "GPLEV")
}

func (chip *Chip2837) twoBankCommand(bcm int, commandName string) (registerAddress uint64, operation int) {
	addressOffset := bcm / 32 //1 register for 32 pins
	shift := (uint8(bcm) % 32)
	operation = 1 << shift
	return chip.gpioRegisters[commandName][addressOffset], operation
}
