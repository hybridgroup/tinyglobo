package main

import (
	"machine"

	"errors"

	"tinygo.org/x/drivers/lora"
	"tinygo.org/x/drivers/sx126x"
)

var (
	spi                        = machine.SPI1
	nssPin, busyPin, dio1Pin   = machine.GP13, machine.GP7, machine.GP6
	rxPin, txLowPin, txHighPin = machine.GP9, machine.GP8, machine.GP8

	loraRadio *sx126x.Device
)

// do sx126x setup here
func setupLora() (lora.Radio, error) {
	spi.Configure(machine.SPIConfig{
		Mode:      0,
		Frequency: 8 * 1e6,
		SDO:       machine.SPI1_SDO_PIN,
		SDI:       machine.SPI1_SDI_PIN,
		SCK:       machine.SPI1_SCK_PIN,
	})

	loraRadio = sx126x.New(spi)
	loraRadio.SetDeviceType(sx126x.DEVICE_TYPE_SX1262)

	loraRadio.SetRadioController(sx126x.NewRadioControl(nssPin, busyPin, dio1Pin, rxPin, txLowPin, txHighPin))

	if state := loraRadio.DetectDevice(); !state {
		return nil, errors.New("LoRa radio not found")
	}

	return loraRadio, nil
}
