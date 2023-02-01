package main

import (
	"machine"

	"tinygo.org/x/drivers/bmp180"
	"tinygo.org/x/drivers/mpu6050"
)

var (
	sensor bmp180.Device
	accel mpu6050.Device

	temperature int32
	pressure int32
	ax, ay, az int32
)

func startSensors() {
	machine.I2C0.Configure(machine.I2CConfig{})

	sensor = bmp180.New(machine.I2C0)
	sensor.Configure()
	accel = mpu6050.New(machine.I2C0)
	accel.Configure()
}

func readSensors() {
	temperature, _ = sensor.ReadTemperature()
	pressure, _ = sensor.ReadPressure()
	ax, ay, az = accel.ReadAcceleration()

	println("Temperature:", float32(temperature)/1000, "°C")
	println("Pressure", float32(pressure)/100000, "hPa")
	println("Accelerometer", ax, ay, az)
}
