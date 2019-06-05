package gopherberry

import (
	"os"
	"syscall"
)

type Epoll struct {
	file  *os.File
	epfd  int
	event syscall.EpollEvent
}

//
// /dev/kmsg
// date | sudo tee /dev/kmsg
//

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
		Events: syscall.EPOLLIN | syscall.EPOLLPRI,
		//Fd:     int32(file.Fd()),
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

func (ep *Epoll) Start() chan struct{} {
	c := make(chan struct{})

	go func() {
		for {
			_, err := syscall.EpollWait(ep.epfd, []syscall.EpollEvent{ep.event}, -1)
			if err != nil {
				//do smth?
			}

			var buf [1024]byte
			_, err = syscall.Read(int(ep.event.Fd), buf[:]) //????? ep.event.Fd -> ep.file.Fd()

			if err != nil {
				//do smth
			}

			c <- struct{}{}
			//fmt.Println(data, e)
		}
	}()
	return c
}

//
func (ep *Epoll) Stop() {
	ep.file.Close()
	syscall.Close(ep.epfd)
}
