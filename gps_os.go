//go:build !tinygo

package main

func gpsIsStarted() bool {
	return false
}

func gpsHasFix() bool {
	return false
}

func initGPS() {
}

// start GPS reading goroutine
func startGPS() {
}

func stopGPS() {
}
