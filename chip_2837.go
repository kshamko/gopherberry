package gopherberry

//Chip2837 implementation for raspberry 3+
type Chip2837 struct {
	//MemMap
	//memMap *Mmap
	//PerBaseAddrPhys - periphial base physical address
	//PerBaseAddrPhys uint64
	//PerBaseAddrBirt - periphial base virtual address
	perBaseAddrVirt uint64
	//Board2BCM maps board pin number to "Broadcom SOC channel" number
	//https://raspberrypi.stackexchange.com/questions/12966/what-is-the-difference-between-board-and-bcm-for-gpio-pin-numbering
	board2BCM map[int]int
	//GPIORegisters maps function to registers
	gpioRegisters map[string][]uint64
}

//NewChip2837 func
func newChip2837() chip {

	//Peripherals (at physical address 0x3F000000 on) are mapped into the kernel virtual address
	//space starting at address 0xF2000000. Thus a peripheral advertised here at bus address
	//0x7Ennnnnn is available in the ARM kenel at virtual address 0xF2nnnnnn.
	chip := &Chip2837{
		//memMap:          mMap,
		//PerBaseAddrPhys: 0x3F000000,
		perBaseAddrVirt: 0xF2000000,
		//addressIncrement:
		board2BCM: map[int]int{
			1:  NoBCMNnm, //3v3 power
			2:  NoBCMNnm, //5v power
			3:  2,        //SDA
			4:  NoBCMNnm, //5v power
			5:  3,        //SCL
			6:  0,        //ground
			7:  4,        //GPCLK0
			8:  14,       //TXD
			9:  NoBCMNnm, //ground
			10: 15,       //RXD
			11: 17,
			12: 18, //PWM0
			13: 27,
			14: NoBCMNnm, //ground
			15: 22,
			16: 23,
			17: NoBCMNnm, //3v3 power
			18: 24,
			19: 10,       //MOSI
			20: NoBCMNnm, //ground
			21: 9,        //MISO
			22: 25,
			23: 11,       //SCLK
			24: 8,        //CE0
			25: NoBCMNnm, //ground
			26: 7,        //CE1
			27: 0,        //ID_SD
			28: 1,        //ID_SC
			29: 5,
			30: NoBCMNnm, //ground
			31: 6,
			32: 12, //PWM0
			33: 13, //PWM1
			34: 16,
			35: 19, //MISO
			36: 16,
			37: 26,
			38: 20,       //MOSI
			39: NoBCMNnm, //ground
			40: 21,       //SCLK
		},
		gpioRegisters: map[string][]uint64{
			"GPFSEL": []uint64{0x7E200000, 0x7E200004, 0x7E200008, 0x7E20000C, 0x7E200010, 0x7E200014},
		},
	}

	return chip
}

func (chip *Chip2837) getPinBCM(pinNumBoard int) int {
	bcm, ok := chip.board2BCM[pinNumBoard]
	if ok {
		return bcm
	}

	return NoBCMNnm
}

func (chip *Chip2837) gpgsel(bcm int, mode pinMode) (addressOffset int, operation int) {

	//calculate proper register offset
	registerOffset := bcm / 10 //1 register for 10 pins

	//calculate operation. all operations are assumed to be 32-bit
	shift := (uint8(bcm) % 10) * 3 // 10 pins per register, command of 3 bits
	operation = int(mode) << shift

	address := chip.gpioRegisters["GPFSEL"][registerOffset]
	return addressOffset, operation
}
