package gopherberry

import (
	"fmt"
	"os/exec"
)

//ModeInput sets pin to input mode
func (p *Pin) ModeInput() error {
	return p.mode(PinModeInput)
}

//DetectEdge func
func (p *Pin) DetectEdge(edge EdgeType) (<-chan EdgeType, error) {
	if p.curMode != PinModeInput {
		return nil, ErrBadPinMode
	}
	p.mu.Lock()
	defer p.mu.Unlock()

	command := fmt.Sprintf("gpio edge %d %s", p.bcmNum, edge)
	_, err := exec.Command("/bin/bash", "-c", command).Output()
	if err != nil {
		return nil, err
	}

	if edge == EdgeNone {
		return nil, p.detectEdgeStop()
	}

	fileName := fmt.Sprintf("/sys/class/gpio/gpio%d/value", p.bcmNum)
	ep, err := NewEpoll(fileName)
	if err != nil {
		return nil, err
	}

	p.epoll = ep
	p.edgeChan = make(chan EdgeType)

	go func() {
		for {
			c := ep.Wait(SeekSet)
			data, ok := <-c
			if ok {

				if data[0] == 49 && (edge == EdgeBoth || edge == EdgeHigh) { //check 1
					p.edgeChan <- EdgeHigh
				}

				if data[0] == 48 && (edge == EdgeBoth || edge == EdgeLow) { //check 0
					p.edgeChan <- EdgeLow
				}
			} else {
				return
			}
		}
	}()

	return p.edgeChan, nil
}

//DetectEdgeStop stop
func (p *Pin) DetectEdgeStop() error {
	if p.curMode != PinModeInput {
		return ErrBadPinMode
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	_, err := p.DetectEdge(EdgeNone)

	return err
}
