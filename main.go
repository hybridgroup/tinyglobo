package main

import (
	"machine"

	"time"

	"tinygo.org/x/drivers/lora/lorawan"
)

var (
	spi                        = machine.SPI1
	nssPin, busyPin, dio1Pin   = machine.GP13, machine.GP7, machine.GP6
	rxPin, txLowPin, txHighPin = machine.GP9, machine.GP8, machine.GP8
)

func main() {
	time.Sleep(5 * time.Second)
	println("*** Sparkie 1 starting... ***")

	// setup LoRa radio
	var err error
	radio, err = setupLora()
	if err != nil {
		failMessage(err)
	}

	// Connect LoRaWAN to use the LoRa Radio device.
	lorawan.UseRadio(radio)

	// Try to connect to the LoRaWAN network
	if err := loraConnect(); err != nil {
		failMessage(err)
	}

	go startGPS()

	startBattery()
	startSensors()

	for {
		println("Sleeping for", LORAWAN_UPLINK_DELAY_SEC, "seconds")
		time.Sleep(time.Second * LORAWAN_UPLINK_DELAY_SEC)

		readBattery()
		readSensors()

		payload, err := createPayload()
		if err != nil {
			println("Payload error:", err)
			continue
		}

		if err := lorawan.SendUplink(payload, session); err != nil {
			println("Uplink error:", err)
			continue
		}
			
		println("Uplink complete, msglen=", len(payload))
	}
}

func failMessage(err error) {
	for {
		println("FATAL:", err)
		time.Sleep(time.Second)
	}
}
