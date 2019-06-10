package gopherberry

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"syscall"

	"golang.org/x/sync/errgroup"
)

var (
	errGraceStop = errors.New("wait grace stop")
)

//Epoll entity
type Epoll struct {
	file     *os.File
	epfd     int
	event    syscall.EpollEvent
	stopChan chan struct{}
	firstRun bool
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
func (ep *Epoll) Wait(ctx1 context.Context) chan []byte {

	c := make(chan []byte)
	ep.stopChan = make(chan struct{})

	ctx, cancel := context.WithCancel(context.Background())
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		<-ep.stopChan
		fmt.Println("return err", errGraceStop)
		cancel()
		return errGraceStop
	})
	g.Go(func() error {
		var buf [1024]byte

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				//could be blocked and stop will not work properly. (triggered on the next iteration)
				//@todo try to implement epoll interrupt with signal call
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

				log.Println("EpollWait")
				//https://support.sas.com/documentation/onlinedoc/sasc/doc750/html/lr1/z2031150.htm
				syscall.Seek(int(ep.event.Fd), 0, 0) //
				i, err := syscall.Read(int(ep.event.Fd), buf[:])
				log.Println("EpollWait1")
				if i == -1 {
					log.Println("Read Num -1")
					continue
				}
				// @todo improve handling
				if err != nil {
					log.Println("Read err:", err)
					continue
				}

				//skip initial epoll cycle
				if ep.firstRun {
					log.Println("first time skip")
					ep.firstRun = false
					continue
				}

				c <- buf[:]
				return nil
			}
		}

	})

	go func(c chan []byte) {
		err := g.Wait()
		log.Println("Wait err:", err)
		close(c)
	}(c)

	return c
}

//Stop func.
// Has known issue when stop happens on the next iteration of EpollWait
func (ep *Epoll) Stop() (err error) {
	ep.stopChan <- struct{}{}
	fmt.Println("send stop signal")
	var x [1024]byte
	syscall.Seek(int(ep.event.Fd), 0, 1) //
	n, err := syscall.Read(int(ep.event.Fd), x[:])
	fmt.Println("1111", n, err)
	return nil
	/*err = syscall.EpollCtl(ep.epfd, syscall.EPOLL_CTL_DEL, int(ep.file.Fd()), &ep.event)
	syscall.Close(ep.epfd) //call to trigger error of EpollWait
	ep.file.Close()

	return nil*/
}
