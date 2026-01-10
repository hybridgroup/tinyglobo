package main

import (
	"machine"
)

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
	temperature = machine.ReadTemperature()
	// pressure, _ = barometer.ReadPressure()
	// altitude, _ = barometer.ReadAltitude()

	// ax, ay, az = accel.ReadAcceleration()

	println("Temperature:", temperature/1000, "°C")
	// println("Pressure", float32(pressure)/100000, "hPa")
	// println("Accelerometer", ax, ay, az)
}
