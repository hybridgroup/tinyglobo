//go:build tinygo

package main

import (
	"machine"
	"time"
)

var watchdog string

func initWatchdog() {
	if watchdog != "enabled" {
		return
	}

	machine.Watchdog.Configure(machine.WatchdogConfig{
		TimeoutMillis: 3000,
	})
	machine.Watchdog.Start()
}

func updateWatchdog() {
	if watchdog != "enabled" {
		return
	}

	machine.Watchdog.Update()
}

func watchAndWait(seconds int) {
	for i := 0; i < seconds; i++ {
		if watchdog == "enabled" {
			machine.Watchdog.Update()
		}
		time.Sleep(1 * time.Second)
	}
}
