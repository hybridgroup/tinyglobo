package main

import (
	"machine"

	"tinygo.org/x/drivers/bmp180"
	"tinygo.org/x/drivers/mpu6050"
)

var (
	// devices
	barometer bmp180.Device
	accel mpu6050.Device

	// readings
	temperature int32
	pressure int32
	altitude int32
	ax, ay, az int32
)

func startSensors() {
	machine.I2C0.Configure(machine.I2CConfig{})

	barometer = bmp180.New(machine.I2C0)
	barometer.Configure()

	accel = mpu6050.New(machine.I2C0)
	accel.Configure()
}

func readSensors() {
	temperature, _ = barometer.ReadTemperature()
	pressure, _ = barometer.ReadPressure()
	altitude, _ = barometer.ReadAltitude()

	ax, ay, az = accel.ReadAcceleration()

	println("Temperature:", float32(temperature)/1000, "°C")
	println("Pressure", float32(pressure)/100000, "hPa")
	println("Accelerometer", ax, ay, az)
}
