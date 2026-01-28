//go:build tinygo

package main

import (
	"machine"
)

var (
	// readings
	temperature int32
	altitude    int32
)

func initSensors() {
}

func readSensors() {
	temperature = machine.ReadTemperature()

	log("Temperature:", temperature/1000, "°C")
}
