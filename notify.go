package main

import (
	"machine"
	"time"
)

type StatusType int

const (
	StatusIdle StatusType = iota
	StatusAcquiringFix
	StatusReadyToTransmit
	StatusTransmitting
	StatusError
)

var Status StatusType = StatusIdle

func notify(count int) {
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	for i := 0; i < count; i++ {
		led.High()
		time.Sleep(250 * time.Millisecond)
		led.Low()
		time.Sleep(250 * time.Millisecond)
	}
}

func startNotification(frequency time.Duration) {
	go func() {
		for {
			notify(int(Status) + 1)
			time.Sleep(frequency)
		}
	}()
}
