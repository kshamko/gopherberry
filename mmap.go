package gopherberry

import (
	"os"
	"sort"
	"sync"
	"syscall"
	"unsafe"
)

//Mmap struct
type mmap struct {
	baseAddress int64
	length      int
	data        []byte
	mu          sync.RWMutex
	datap       *[100]int
	offsets     map[uint64]int
}

func newMmap(addressesPhysical []uint64) (*mmap, error) {
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

	baseAddress, length, offsets := mmapParameters(addressesPhysical)

	data, err := syscall.Mmap(
		int(file.Fd()),
		baseAddress,
		length,
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_SHARED,
	)

	mmapArray := (*[100]int)(unsafe.Pointer(&data[0]))

	if err != nil {
		return nil, err
	}

	return &mmap{
		data:        data,
		baseAddress: baseAddress,
		length:      length,
		mu:          sync.RWMutex{},
		datap:       mmapArray,
		offsets:     offsets,
	}, nil
}

func (mmap *mmap) run(address uint64, command int) error {

	mmap.mu.Lock()
	defer mmap.mu.Unlock()

	if offset, ok := mmap.offsets[address]; ok {
		mmap.datap[offset] = command
		return nil
	}

	return ErrNoMmap
}

//
func (mmap *mmap) get(address uint64) (state int, err error) {
	mmap.mu.Lock()
	defer mmap.mu.Unlock()

	if offset, ok := mmap.offsets[address]; ok {
		return mmap.datap[offset], nil
	}

	return 0, ErrNoMmap
}

//
func mmapParameters(addressesPhysical []uint64) (baseAddress int64, length int, offsets map[uint64]int) {
	sort.Slice(addressesPhysical, func(i, j int) bool { return addressesPhysical[i] < addressesPhysical[j] })

	for i, addr := range addressesPhysical {
		offsets[addr] = i
	}
	return int64(addressesPhysical[0]), len(addressesPhysical), offsets
}
