package main

import (
	"fmt"

	"github.com/kshamko/gopherberry"
)

//
// Trigger it with "$ date | sudo tee /dev/kmsg"
//
func main() {

	defer fmt.Println("Stopped")

	ep, _ := gopherberry.NewEpoll("/dev/kmsg")
	c := ep.Start()

	fmt.Println("started")
	x := 0
	for {
		select {
		case _, ok := <-c:

			if !ok {
				fmt.Println("closed")
				return
			}
			x++
			fmt.Println("changed", x, "times")

			if x == 3 {
				ep.Stop()
			}
		}
	}
}

/*import (
	"fmt"
	"os"
	"syscall"
)

const (
	MaxEpollEvents = 32
	KB             = 1024
)

func echo(in, out int) {

	var buf [KB]byte
	for {
		nbytes, e := syscall.Read(in, buf[:])
		if nbytes > 0 {
			syscall.Write(out, buf[:nbytes])
		}
		if e != nil {
			break
		}
	}
}

func main() {

	var event syscall.EpollEvent
	var events [MaxEpollEvents]syscall.EpollEvent

	file, e := os.Open("/dev/kmsg")
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
	defer file.Close()

	fd := int(file.Fd())
	if e = syscall.SetNonblock(fd, true); e != nil {
		fmt.Println("setnonblock1: ", e)
		os.Exit(1)
	}

	epfd, e := syscall.EpollCreate(1)
	if e != nil {
		fmt.Println("epoll_create1: ", e)
		os.Exit(1)
	}
	defer syscall.Close(epfd)

	event.Events = syscall.EPOLLIN | syscall.EPOLLPRI
	event.Fd = int32(fd)
	if e = syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, fd, &event); e != nil {
		fmt.Println("epoll_ctl: ", e)
		os.Exit(1)
	}

	for {

		_, e := syscall.EpollWait(epfd, events[:], -1)
		if e != nil {

			fmt.Println("epoll_wait: ", e)
			os.Exit(1)
			break
		}

		var buf [1024]byte
		data, e := syscall.Read(int(events[0].Fd), buf[:])
		fmt.Println(data, e)

		//echo(int(events[0].Fd), syscall.Stdout)
		/*fmt.Println(nevents)
		for ev := 0; ev < nevents; ev++ {
			go echo(int(events[ev].Fd), syscall.Stdout)
		}

	}
}*/
