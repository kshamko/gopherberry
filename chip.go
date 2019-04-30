package gopherberry

import (
	"errors"
	//"github.com/kshamko/gopherberry/gpio"
	//"os"
)

const (
	//NoBCMNnm means that pin has no bcm number (ground, voltage pins)
	NoBCMNnm = -1
)

var (
	//ErrNoPin error prodices when there is no bcm number for pin
	ErrNoPin = errors.New("no pin exists")
)

//Chip struct
type Chip struct {
	//MemMap
	memMap *Mmap
	//PerBaseAddrPhys - periphial base physical address
	PerBaseAddrPhys uint64
	//PerBaseAddrBirt - periphial base virtual address
	PerBaseAddrVirt uint64
	//Board2BCM maps board pin number to "Broadcom SOC channel" number
	//https://raspberrypi.stackexchange.com/questions/12966/what-is-the-difference-between-board-and-bcm-for-gpio-pin-numbering
	Board2BCM     map[int]int
	GPIORegisters map[string][]uint64
}

//NewChip2837 func
func NewChip2837() (*Chip, error) {

	mMap, err := NewMmap()

	//Thus a peripheral advertised here at bus address 0x7Ennnnnn is available in the ARM kenel at virtual address 0xF2nnnnnn
	chip := &Chip{
		memMap:          mMap,
		PerBaseAddrPhys: 0x3F000000,
		PerBaseAddrVirt: 0xF2000000,
		Board2BCM: map[int]int{
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
		GPIORegisters: map[string][]uint64{},
	}

	return chip, err
}

//GetPin retuns pin object
func (c *Chip) GetPin(pinNumBoard int) (*Pin, error) {

	if num, ok := c.Board2BCM[pinNumBoard]; ok && num != NoBCMNnm {
		return &Pin{
			BCMNum: num,
			mMap:   c.memMap,
		}, nil
	}

	return nil, ErrNoPin
}
