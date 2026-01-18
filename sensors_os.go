//go:build !tinygo

package main

var (
	// readings
	temperature int32
	pressure    int32
	altitude    int32
	ax, ay, az  int32
)

func initSensors() {
}

func readSensors() {
}
