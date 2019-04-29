package gopherberry

import (
	"fmt"
	"os"
)

//ChipMmap struct
type ChipMmap struct {
}

//NewMmap func
func NewMmap() (*ChipMmap, error) {
	//var file *os.File

	// Open fd for rw mem access; try dev/mem first (need root)
	file, err := os.OpenFile("/dev/mem", os.O_RDWR|os.O_SYNC, 0)
	/*if os.IsPermission(err) { // try gpiomem otherwise (some extra functions like clock and pwm setting wont work)
		file, err = os.OpenFile("/dev/gpiomem", os.O_RDWR|os.O_SYNC, 0)
	}*/
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// FD can be closed after memory mapping
	defer file.Close()

	/*memlock.Lock()
	defer memlock.Unlock()
	// Memory map GPIO registers to slice
	gpioMem, gpioMem8, err = memMap(file.Fd(), gpioBase)
	if err != nil {
		return
	}
	// Memory map clock registers to slice
	clkMem, clkMem8, err = memMap(file.Fd(), clkBase)
	if err != nil {
		return
	}
	// Memory map pwm registers to slice
	pwmMem, pwmMem8, err = memMap(file.Fd(), pwmBase)
	if err != nil {
		return
	}
	// Memory map spi registers to slice
	spiMem, spiMem8, err = memMap(file.Fd(), spiBase)
	if err != nil {
		return
	}
	// Memory map interruption registers to slice
	intrMem, intrMem8, err = memMap(file.Fd(), intrBase)
	if err != nil {
		return
	}*/
	//backupIRQs() // back up enabled IRQs, to restore it later

	return nil, nil
}
