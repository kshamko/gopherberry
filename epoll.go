package gopherberry

import (
	"log"
	"os"
	"syscall"
)

//FilePos to tell read system call where to start reading
type FilePos int

const (
	//SeekSet is the beginning of the file; the value is 0.
	SeekSet FilePos = iota
	//SeekCur is the current file offset; the value is 1.
	SeekCur
	//SeekEnd is the end of the file; the value is 2.
	SeekEnd
)

//Epoll entity
type Epoll struct {
	file     *os.File
	epfd     int
	firstRun bool
	event    syscall.EpollEvent
}

//NewEpoll func
func NewEpoll(fileName string) (*Epoll, error) {

	//https://medium.com/coderscorner/tale-of-client-server-and-socket-a6ef54a74763
	file, err := os.OpenFile(fileName, syscall.O_RDONLY|syscall.O_NONBLOCK, 0644)
	if err != nil {
		return nil, err
	}

	epfd, err := syscall.EpollCreate1(0)
	if err != nil {
		return nil, err
	}

	EPOLLET := uint32(1 << 31)
	event := syscall.EpollEvent{
		Events: syscall.EPOLLIN | syscall.EPOLLPRI | EPOLLET | syscall.EPOLLERR,
		Fd:     int32(file.Fd()),
	}

	err = syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, int(file.Fd()), &event)
	if err != nil {
		return nil, err
	}

	return &Epoll{
		file:     file,
		epfd:     epfd,
		event:    event,
		firstRun: true,
	}, nil
}

//Wait func
func (ep *Epoll) Wait(pos FilePos) chan []byte {

	c := make(chan []byte)

	go func() {
		var buf [1024]byte

		for {
			num, err := syscall.EpollWait(ep.epfd, []syscall.EpollEvent{ep.event}, -1)
			if num == -1 {
				log.Println("EpollWait Num -1")
				continue
			}
			// @todo improve handling
			if err != nil {
				log.Println("EpollWait err:", err)
				continue
			}

			//https://support.sas.com/documentation/onlinedoc/sasc/doc750/html/lr1/z2031150.htm
			seek(int(ep.event.Fd), int(pos), ep.firstRun)

			i, err := syscall.Read(int(ep.event.Fd), buf[:])
			if i == -1 {
				log.Println("Read Num -1")
				continue
			}
			// @todo improve handling
			if err != nil {
				log.Println("Read err:", err)
				continue
			}

			//skip the first epoll cycle for pos types not equal to SeekEnd
			if ep.firstRun && pos != SeekEnd {
				ep.firstRun = false
				continue
			}

			ep.firstRun = false
			c <- buf[:]

			return
		}
	}()

	return c
}

//Stop func.
func (ep *Epoll) Stop() error {
	err := syscall.EpollCtl(ep.epfd, syscall.EPOLL_CTL_DEL, int(ep.file.Fd()), &ep.event)
	_ = syscall.Close(ep.epfd)
	_ = ep.file.Close()
	return err
}

//seek systemcall wrapper
func seek(fd int, pos int, firstRun bool) {
	if pos == int(SeekEnd) {
		if !firstRun {
			return
		}
		_, _ = syscall.Seek(fd, 0, pos)
		return
	}
	_, _ = syscall.Seek(fd, 0, pos)
}
