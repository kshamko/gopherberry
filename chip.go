package gopherberry

import "github.com/kshamko/gopherberry/gpio"

//Chip struct
type Chip struct {
	//PerBaseAddrPhys - periphial base physical address
	PerBaseAddrPhys int64
	//PerBaseAddrBirt - periphial base virtual address
	PerBaseAddrVirt int64
	//Board2BCM maps board pin number to "Broadcom SOC channel" number
	//https://raspberrypi.stackexchange.com/questions/12966/what-is-the-difference-between-board-and-bcm-for-gpio-pin-numbering
	Board2BCM map[int]int
}

//NewChip func
func NewChip() *Chip {

	//Thus a peripheral advertised here at bus address 0x7Ennnnnn is available in the ARM kenel at virtual address 0xF2nnnnnn
	return &Chip{
		PerBaseAddrPhys: 0x3F000000,
		PerBaseAddrVirt: 0xF2000000,
		Board2BCM: map[int]int{
			1: 0,
			2: 0,
			3: 0,
			//..,
			40: 0,
		},
	}
}

//GetPin retuns pin object
func (c *Chip) GetPin(pinNum int) (*gpio.Pin, error) {

	return nil, nil
}
