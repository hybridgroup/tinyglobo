package main

import (
	"machine"

	"errors"

	"tinygo.org/x/drivers/lora"
	"tinygo.org/x/drivers/sx126x"
)

const FREQ                      = 868100000

var (
	loraRadio *sx126x.Device
)

// do sx126x setup here
func SetupLora() (lora.Radio, error) {
	spi.Configure(machine.SPIConfig{
		Mode:      0,
		Frequency: 8 * 1e6,
		SDO:       machine.SPI1_SDO_PIN,
		SDI:       machine.SPI1_SDI_PIN,
		SCK:       machine.SPI1_SCK_PIN,
	})

	loraRadio = sx126x.New(spi)
	loraRadio.SetDeviceType(sx126x.DEVICE_TYPE_SX1262)

	// Create radio controller for target
	loraRadio.SetRadioController(newRadioControl())

	if state := loraRadio.DetectDevice(); !state {
		return nil, errors.New("LoRa radio not found")
	}

	loraConf := lora.Config{
		Freq:           FREQ,
		Bw:             lora.Bandwidth_125_0,
		Sf:             lora.SpreadingFactor9,
		Cr:             lora.CodingRate4_7,
		HeaderType:     lora.HeaderExplicit,
		Preamble:       12,
		Iq:             lora.IQStandard,
		Crc:            lora.CRCOn,
		SyncWord:       lora.SyncPublic,
		LoraTxPowerDBm: 20,
	}

	loraRadio.LoraConfig(loraConf)

	return loraRadio, nil
}

func newRadioControl() sx126x.RadioController {
	return sx126x.NewRadioControl(nssPin, busyPin, dio1Pin, rxPin, txLowPin, txHighPin)
}
