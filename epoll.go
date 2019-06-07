package gopherberry

import (
	"fmt"
	"os"
	"syscall"
)

//Epoll entity
type Epoll struct {
	file  *os.File
	epfd  int
	event syscall.EpollEvent
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
		file:  file,
		epfd:  epfd,
		event: event,
	}, nil
}

//Wait func
func (ep *Epoll) Wait() chan []byte {

	c := make(chan []byte)

	go func() {
		var buf [1]byte
		for {
			//could be blocked and stop will not work properly. (on the next iteration)
			//@todo try to implement epoll interrupt with signal call
			num, err := syscall.EpollWait(ep.epfd, []syscall.EpollEvent{ep.event}, -1)

			fmt.Println("!!!!i!!!", err, num)
			if num == -1 {
				//continue
				return
			}
			// @todo improve handling
			if err != nil {
				close(c)
				//ep.Stop()

				return
			}
			//
			//https://support.sas.com/documentation/onlinedoc/sasc/doc750/html/lr1/z2031150.htm
			syscall.Seek(int(ep.event.Fd), 0, 2)
			i, err := syscall.Read(int(ep.event.Fd), buf[:])

			if i == -1 {
				fmt.Println("!!!!i!!!", i)
				continue
				//return
			}
			if err != nil {
				fmt.Println("!!!!i!!!", err)
				close(c)
				//ep.Stop()

				return
				//do smth
			}
			c <- buf[:]
			close(c)
			return
			//ep.Stop()
			//return
		}
	}()

	return c
}

//Stop func.
// Has known issue when stop happens on the next iteration of EpollWait
func (ep *Epoll) Stop() error {
	syscall.Close(ep.epfd) //call to trigger error of EpollWait
	ep.file.Close()
	syscall.EpollCtl(ep.epfd, syscall.EPOLL_CTL_DEL, int(ep.file.Fd()), &ep.event)
	return nil
}
