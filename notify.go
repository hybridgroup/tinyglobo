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

func notifyLowBattery() {
	stopNotifications()
	notify(int(StatusLowBattery))
	time.Sleep(500 * time.Millisecond)
	notify((int(currentVoltage) / 1000) - 1)
	time.Sleep(500 * time.Millisecond)
}

func notifyError(status StatusType) {
	stopNotifications()
	notify(int(status))
	time.Sleep(500 * time.Millisecond)
	notify(int(status))
	time.Sleep(500 * time.Millisecond)
	notify(int(status))
	time.Sleep(500 * time.Millisecond)
}
