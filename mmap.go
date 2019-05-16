package gopherberry

import (
	//"bytes"

	"fmt"
	"syscall"

	//"launchpad.net/gommap"
	"os"
)

//Mmap struct
type mmap struct {
	memMap map[string][]byte
}

//NewMmap func
func newMmap(registers map[string][]uint64, baseVirtAddress uint64) (*mmap, error) {

	// Open fd for rw mem access; try dev/mem first (need root)
	file, err := os.OpenFile("/dev/kmem", os.O_RDWR|os.O_SYNC, 0)
	/*if os.IsPermission(err) { // try gpiomem otherwise (some extra functions like clock and pwm setting wont work)
		file, err = os.OpenFile("/dev/gpiomem", os.O_RDWR|os.O_SYNC, 0)
	}*/

	if err != nil {
		return nil, err
	}

	// FD can be closed after memory mapping
	defer file.Close()

	memMap := map[string][]byte{}
	for op, addresses := range registers {

		fmt.Println("[INFO] ", op, "address:", virtAddress(addresses[0], baseVirtAddress), "div: ", virtAddress(addresses[0], baseVirtAddress)%uint64(os.Getpagesize()))

		data, err := mapMemory(file.Fd(), virtAddress(addresses[0], baseVirtAddress), len(addresses))
		if err == nil {
			memMap[op] = data
		}
	}

	return &mmap{memMap}, nil
}

//
func (mm *mmap) run(funcName string, registerOffset int, command int) error {
	mm.memMap[funcName][registerOffset] = byte(command)
	return nil
}

//
func mapMemory(fd uintptr, base uint64, len int) ([]byte, error) {
	return syscall.Mmap(
		int(fd),
		int64(base),
		len,
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_SHARED|syscall.MAP_FIXED,
	)
}

func virtAddress(busAddress uint64, baseVirtAddress uint64) uint64 {
	base := busAddress & 0xff000000
	return baseVirtAddress + (busAddress - base)
}
