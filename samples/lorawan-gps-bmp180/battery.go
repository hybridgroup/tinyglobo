package main

import (
	"machine"
)

var (
	vsys = machine.ADC{machine.ADC3}
	voltage uint32
)

func startBattery() {
	machine.InitADC()	
	vsys.Configure(machine.ADCConfig{})
}

func readBattery() {
	// calculate in millivolts
	voltage = uint32(vsys.Get())*10*323*3/65535
}
