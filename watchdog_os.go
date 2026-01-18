//go:build !tinygo

package main

import (
	"time"
)

var watchdog string

func initWatchdog() {
}

func updateWatchdog() {
}

func watchAndWait(seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
}
