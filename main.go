package main

import (
	"time"

	"tinygo.org/x/wireless/wspr"
)

func main() {
	time.Sleep(5 * time.Second)
	println("*** TinyGlobo 3 starting... ***")

	go startGPS()

	// setup radio
	var err error
	radio, err := initRadio()
	if err != nil {
		failMessage(err)
	}

	frequency := radio.GetBaseFrequency()
	println("Transmitting on frequency", frequency, "Hz")

	data := make([]byte, 162)

	startBattery()
	startSensors()

	for {
		// println("Sleeping for", uplinkDelaySeconds, "seconds")
		// time.Sleep(time.Second * uplinkDelaySeconds)
		println("Waiting for next transmission...")
		time.Sleep(15 * time.Second)

		readBattery()
		readSensors()

		// Example WSPR packet data
		// K1ABC FN42 37
		// See https://en.wikipedia.org/wiki/WSPR_(amateur_radio_software)
		msg, err := wspr.NewMessage("K1ABC", "FN42", 37)
		if err != nil {
			println("Error creating WSPR message:", err.Error())
			return
		}

		n, err := msg.WriteSymbols(data)
		if err != nil {
			println("error writing WSPR message")
			return
		}

		println("Transmitting WSPR message with", n, "symbols")
		if err := radio.WriteSymbols(data[:n]); err != nil {
			println("error transmitting WSPR message:", err.Error())
			return
		}
	}
}

func failMessage(err error) {
	for {
		println("FATAL:", err)
		time.Sleep(time.Second)
	}
}
