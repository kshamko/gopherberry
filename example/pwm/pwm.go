package main

import (
	"fmt"
	"time"

	"github.com/kshamko/gopherberry"
)

func main() {

	r, err := gopherberry.New(gopherberry.ARM2837)

	if err != nil {
		fmt.Println("[ERROR] can't init pi", err)
	}

	
	p18, _ := r.GetPin(12)
	err = p18.ModePWM()
	if err != nil {
		fmt.Println("[ERROR] cant set mode to pin 18(12)", err)
	}

	p18.SetFrequency(gopherberry.ClockConfig{Enab: true}, 32000)

	c := gopherberry.PWMChannelConfig{
		MSEnable:    1,
		ChanEnabled: 1,
	}
	err = r.StartPWM(c, gopherberry.PWMChannelConfig{})
	if err != nil {
		fmt.Println("[ERROR] can't init pwm", err)
	}
	

	p18.DutyCycle(32, 32)
	// the LED will be blinking at 2000Hz
	// (source frequency divided by cycle length => 64000/32 = 2000)

	time.Sleep(10 * time.Second)
	// five times smoothly fade in and out
	for i := 0; i < 5; i++ {
		for i := int(0); i < 32; i++ { // increasing brightness
			p18.DutyCycle(i, 32)
			time.Sleep(time.Second / 32)
		}
		for i := int(32); i > 0; i-- { // decreasing brightness
			p18.DutyCycle(i, 32)
			time.Sleep(time.Second / 32)
		}
	}

}
