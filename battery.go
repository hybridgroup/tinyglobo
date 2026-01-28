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
	var voltage uint16
	// for i := 0; i < 10; i++ {
	// 	v := vsys.Get()
	// 	if v > voltage {
	// 		voltage = v
	// 	}
	// 	time.Sleep(5 * time.Millisecond)
	// }

	voltage = vsys.Get()
	currentVoltage = uint32(voltage) * 10 * 323 * 3 / 65535

	return currentVoltage
}
