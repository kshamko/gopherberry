package gopherberry

import (
	"os"
	"sync"
	"syscall"
	"fmt"
	"unsafe"
)

//Mmap struct
type mmap struct {
	baseAddress int64
	length      int
	data        []byte
	mu          sync.RWMutex
	datap       *[200]int
}

func newMmap(baseAddress int64, length int) (*mmap, error) {
	// to open /dev/mem we need root. as we use /dev/mem => we'll use physical addresses
	// http://man7.org/linux/man-pages/man4/mem.4.html
	file, err := os.OpenFile("/dev/mem", os.O_RDWR|os.O_SYNC, 0)
	if os.IsPermission(err) { // if we have no root
		file, err = os.OpenFile("/dev/gpiomem", os.O_RDWR|os.O_SYNC, 0)
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := syscall.Mmap(
		int(file.Fd()),
		baseAddress,
		length, //os.Getpagesize(),
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_SHARED,
	)

fmt.Println("!!!!", length, data)
	
			   map_array := (*[200]int)(unsafe.Pointer(&data[0]))

			    /*for i := 0; i < n; i++ {
			        map_array[i] = i * i
				}

				//after all

		    err = syscall.Munmap(data)
		    if err != nil {
		        fmt.Println(err)
		        os.Exit(1)
		    }
	*/

	if err != nil {
		return nil, err
	}

	return &mmap{
		data:        data,
		baseAddress: baseAddress,
		length:      length,
		mu:          sync.RWMutex{},
		datap: map_array,
	}, nil
}

func (mmap *mmap) run(offset int, command int) error {
	if len(mmap.data) < offset {
		return ErrNoMmap
	}

	mmap.mu.Lock()
	mmap.datap[offset] = int(mmap.datap[offset]) ^ command
	mmap.mu.Unlock()

	return nil
}
