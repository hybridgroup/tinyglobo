package main

import (
	"machine"
)

const (
	minBatteryVoltage     = 2500 // millivolts
	desiredBatteryVoltage = 3000 // millivolts
)

var (
	vsys = machine.ADC{machine.ADC3}

	// battery voltage in millivolts
	voltage uint32
)

func initBattery() {
	vsys.Configure(machine.ADCConfig{})
}

func readBattery() {
	// calculate in millivolts
	voltage = uint32(vsys.Get()) * 10 * 323 * 3 / 65535
}
