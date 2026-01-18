//go:build tinygo

package main

import (
	"machine"
)

var (
	vsys = machine.ADC{machine.ADC3}
)

func initBattery() {
	machine.InitADC()

	vsys.Configure(machine.ADCConfig{})
}

// calculate in millivolts
func readBattery() uint32 {
	currentVoltage = uint32(vsys.Get()) * 10 * 323 * 3 / 65535

	return currentVoltage
}
