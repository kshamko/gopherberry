package gopherberry

import (
	"fmt"
	"os"
	"syscall"
)

const (
	EPOLLERR = 0x008
)

//Epoll entity
type Epoll struct {
	file     *os.File
	epfd     int
	event    syscall.EpollEvent
	stopChan chan struct{}
}

//NewEpoll func
func NewEpoll(fileName string) (*Epoll, error) {

	//https://medium.com/coderscorner/tale-of-client-server-and-socket-a6ef54a74763
	file, err := os.OpenFile(fileName, syscall.O_RDONLY|syscall.O_NONBLOCK, 0644)
	if err != nil {
		return nil, err
	}

	epfd, err := syscall.EpollCreate(1)
	if err != nil {
		return nil, err
	}

	event := syscall.EpollEvent{
		Events: syscall.EPOLLIN | syscall.EPOLLPRI | EPOLLERR,
		Fd:     int32(file.Fd()),
	}

	err = syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, int(file.Fd()), &event)
	if err != nil {
		return nil, err
	}

	return &Epoll{
		file:  file,
		epfd:  epfd,
		event: event,
	}, nil
}

//Start func
func (ep *Epoll) Start() chan []byte {

	c := make(chan []byte)
	ep.stopChan = make(chan struct{}, 1)
	syscall.Seek(int(ep.event.Fd), 0, 2)

	go func() {
		var buf [1024]byte
		for {

			select {
			case <-ep.stopChan:
				syscall.EpollCtl(ep.epfd, syscall.EPOLL_CTL_DEL, int(ep.file.Fd()), &ep.event)
				ep.file.Close()
				close(c)
				return
			default:
				//could be blocked and stop will not work properly. (on the next iteration)
				//@todo try to implement epoll interrupt with signal call
				num, err := syscall.EpollWait(ep.epfd, []syscall.EpollEvent{ep.event}, -1)

				if num == -1 {
					continue
				}
				// @todo improve handling
				if err != nil {
					continue
					//do smth?
				}
				//
				i, err := syscall.Read(int(ep.event.Fd), buf[:])
				if i == -1 {
					continue
				}
				if err != nil {
					fmt.Println("err:", err)
					continue
					//do smth
				}
				c <- buf[:]
			}
		}
	}()

	return c
}

//Stop func.
// Has known issue when stop happens on the next iteration of EpollWait
func (ep *Epoll) Stop() error {

	err := syscall.Close(ep.epfd) //call to trigger error of EpollWait
	ep.stopChan <- struct{}{}

	return err
}
