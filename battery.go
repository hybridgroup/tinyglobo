//go:build tinygo

package main

import (
	"machine"
)

var (
	vsys = machine.ADC{machine.ADC3}

	// battery voltage in millivolts
	voltage uint32
)

func initBattery() {
	machine.InitADC()

	vsys.Configure(machine.ADCConfig{})
}

// calculate in millivolts
func readBattery() uint32 {
	voltage = uint32(vsys.Get()) * 10 * 323 * 3 / 65535

	return voltage
}
