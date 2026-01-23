//go:build tinygo

package main

import (
	"machine"
	"time"
)

func notify(status int) {
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	for i := 0; i < status+1; i++ {
		led.High()
		time.Sleep(250 * time.Millisecond)
		led.Low()
		time.Sleep(250 * time.Millisecond)
	}
}
