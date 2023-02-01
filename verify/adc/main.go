package main

import (
	"machine"
	"time"
)

func main() {
	machine.InitADC()

	voltage := machine.ADC{machine.ADC3}
	voltage.Configure(machine.ADCConfig{})

	for {
		// calculate in millivolts
		vsys := uint32(voltage.Get())*10*323*3/65535
		println(vsys)
		time.Sleep(time.Millisecond * 500)
	}
}
