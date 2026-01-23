package main

import (
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

func startNotification(frequency time.Duration) {
	go func() {
		for {
			notify(int(Status))
			time.Sleep(frequency)
		}
	}()
}
