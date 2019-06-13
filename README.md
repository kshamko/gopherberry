# Raspberry Pi Programming With Go

Golang SDK for Raspberry Pi programming. **Mmap** <del>(it's not a complete truth)</del> approach was used so it's important to read the spec of BCM chip, which controls Pi. Spec could be found [here](https://github.com/kshamko/gopherberry/BCM2837-ARM-Peripherals.-.Revised.-.V2-1.pdf)

## Description

Each GPIO pin supports up to 8 modes. The full description of the modes could be found in the spec on page *102*. At the momend SDK supports 2 common modes for each pin: **input** & **output**. On lower level pins have a set functions which could be activated with memory mapped registers. 
Refer to the spec's page *90* for a complete list of functions and registers' addresses.

To be able to make an edge detection on pins mmap approach was declined because there is no way to catch a CPU interrupt from Go program and gracefully process it. With mmap there is only solution to disable interrupts (page 112) to prevent Pi from hanging and poll *GPREN* register in an infinite loop. So **epoll** approach was used to detect rising/falling edge on a pin. More details could be found in the presentation:

```bash
$ cd presentation
$ present
```

## Runnig Gopherberry

It's important to know that pins have a modes:
- to set level of pins from Go program the output mode should be used
- to detect level on pins please use input mode (i.e. use to detect signals from sensors) 
- each pin has it's own number of 2 types: *board number* & *bcm number*. [details](https://pinout.xyz/#). **This SDK uses board number to initialise a pin**
  ![Board Vs BCM Num](/docs/pins.png)
- more details about pins could be found [here](https://pinout.xyz/#)

Please refer to the examples for running tips



