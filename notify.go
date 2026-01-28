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
	StatusLowBattery
	StatusGeofenceBreach
	StatusTooSoonToTransmit
	StatusError
)

var Status StatusType = StatusIdle

var notifyStopChan chan struct{}

func startNotifications(frequency time.Duration) {
	notifyStopChan = make(chan struct{})

	go func() {
		for {
			select {
			case <-notifyStopChan:
				return
			default:
			}

			notify(int(Status))
			time.Sleep(frequency)
		}
	}()
}

func stopNotifications() {
	if notifyStopChan != nil {
		close(notifyStopChan)
		notifyStopChan = nil
	}
}
