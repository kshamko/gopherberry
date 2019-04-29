package gopherberry

import (
	//"fmt"
	"github.com/kshamko/gopherberry/gpio"
	//"os"
)

//Chip struct
type Chip struct {
	//MemMap
	MemMap *ChipMmap
	//PerBaseAddrPhys - periphial base physical address
	PerBaseAddrPhys int64
	//PerBaseAddrBirt - periphial base virtual address
	PerBaseAddrVirt int64
	//Board2BCM maps board pin number to "Broadcom SOC channel" number
	//https://raspberrypi.stackexchange.com/questions/12966/what-is-the-difference-between-board-and-bcm-for-gpio-pin-numbering
	Board2BCM map[int]int
}

//NewChip func
func NewChip() (*Chip, error) {

	mMap, err := NewMmap()

	//Thus a peripheral advertised here at bus address 0x7Ennnnnn is available in the ARM kenel at virtual address 0xF2nnnnnn
	chip := &Chip{
		MemMap:          mMap,
		PerBaseAddrPhys: 0x3F000000,
		PerBaseAddrVirt: 0xF2000000,
		Board2BCM: map[int]int{
			1:  -1, //3v3 power
			2:  -1, //5v power
			3:  2,  //SDA
			4:  -1, //5v power
			5:  3,  //SCL
			6:  0,  //ground
			7:  4,  //GPCLK0
			8:  14, //TXD
			9:  -1, //ground
			10: 15, //RXD
			11: 17,
			12: 18, //PWM0
			13: 27,
			14: -1, //ground
			15: 22,
			16: 23,
			17: -1, //3v3 power
			18: 24,
			19: 10, //MOSI
			20: -1, //ground
			21: 9,  //MISO
			22: 25,
			23: 11, //SCLK
			24: 8,  //CE0
			25: -1, //ground
			26: 7,  //CE1
			27: 0,  //ID_SD
			28: 1,  //ID_SC
			29: 5,
			30: -1, //ground
			31: 6,
			32: 12, //PWM0
			33: 13, //PWM1
			34: 16,
			35: 19, //MISO
			36: 16,
			37: 26,
			38: 20, //MOSI
			39: -1, //ground
			40: 21, //SCLK
		},
	}

	return chip, err
}

//GetPin retuns pin object
func (c *Chip) GetPin(pinNum int) (*gpio.Pin, error) {

	return nil, nil
}
