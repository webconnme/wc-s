package main

import (
	"wc/gpio"
	"log"
)

type LedInfo struct {
	module uint
	bit uint
}

var leds = []LedInfo  {
	{1, 31},
	{2, 0},
	{2, 1},
}

func main() {
	g, err := gpio.Open()
	if err != nil {
		log.Panic("GpioInit: ", err)
	}

	for _, led := range leds {
		g.Alt(led.module, led.bit, 0)
		g.Dir(led.module, led.bit, 1)
	}

	for _, led := range leds {
		g.Val(led.module, led.bit, 1)
	}

	err = g.Close()
	if err != nil {
		log.Panic("GpioDeinit: ", err)
	}
}
