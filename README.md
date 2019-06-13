# Raspberry Pi Programming With Go

Golang SDK for Raspberry Pi programming. **Mmap** <del>(it's not a complete truth)</del> approach was used so it's important to read the spec of BCM chip, which controls Pi. Spec could be found [here](/docs/BCM2837-ARM-Peripherals.-.Revised.-.V2-1.pdf)

## Description

Each GPIO pin supports up to 8 modes. The full description of the modes could be found in the spec on page *102*. At the moment SDK supports 2 common modes for each pin: **input** & **output**. On lower level pins have a set functions which could be activated with memory mapped registers. 
Refer to the spec's page *90* for a complete list of functions and registers' addresses. Spec mentions bus addresses which should be recalculated into physical ones (because memory mapped */dev/mem* file is about physical addresses). This could be done using the following rule from the spec:

> Peripherals (at physical address 0x3F000000 on) are mapped into the kernel virtual address space starting at address 0xF2000000. Thus a peripheral advertised here at bus address 0x7Ennnnnn is available in the ARM kernel at virtual address 0xF2nnnnnn.

To be able to make an edge detection on pins mmap approach was declined because there is no way to catch a CPU interrupt from Go program and gracefully process it. With mmap there is the only solution to disable interrupts (page 112) to prevent Pi from hanging and poll *GPREN* register in an infinite loop. So **epoll** approach was used to detect rising/falling edge on a pin. More details could be found in the presentation:

```bash
$ cd presentation
$ present
```

## Runnig Gopherberry

It's important to know that pins have a modes:
- to set level of pins from Go program the output mode should be used
- to detect level on pins please use input mode (i.e. use to detect signals from sensors) 
- each pin has it's own number of 2 types: *board number* & *bcm number*.
**This SDK uses board number to initialise a pin**

  ![Board Vs BCM Num](/docs/pins.png)
  
- more details about pins could be found [here](https://pinout.xyz/#)

Please refer to the examples for running tips



