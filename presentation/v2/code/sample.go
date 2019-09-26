package code

import (
	"os"
	"syscall"
	"unsafe"
)

func loadMmaped(baseAddress int64) ([]byte, error) {
	file, err := os.OpenFile("/dev/mem", os.O_RDWR|os.O_SYNC, 0)
	if os.IsPermission(err) { // if we have no root
		file, err = os.OpenFile("/dev/gpiomem", os.O_RDWR|os.O_SYNC, 0)
	}

	data, err := syscall.Mmap(
		int(file.Fd()),
		baseAddress,
		os.Getpagesize(),
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_SHARED,
	)

	return data, err
}

func toPointer(data []byte) *[100]int {
	return (*[100]int)(unsafe.Pointer(&data[0]))
}

func writeCommand(mem *[100]int, offset int, command int) {
	mem[offset] = command
}
