//go:build !tinygo

package main

import "time"

func gpsIsStarted() bool {
	return false
}

func gpsHasFix() bool {
	return false
}

func gpsTimeAdjusted() bool {
	return false
}

func initGPS() {
}

// start GPS reading goroutine
func startGPS() {
}

func stopGPS() {
}

func adjustTimeFromGPS(newTime time.Time) {
}
