package main

import (
	"tinygo.org/x/wireless/u4b"
	"tinygo.org/x/wireless/wspr"
)

const dialFreq = 14_095_600 // in Hz

// Determine the frequency band to use for transmission.
// band 3 is skipped to avoid interference.
func frequencyBand(band int) int {
	switch band {
	case 1, 2:
		return band
	case 3, 4:
		return band + 1
	default:
		return 1
	}
}

// transmitFrequency returns the transmit frequency in Hz for the given band.
func transmitFrequency(band int) uint64 {
	freqTxLow := dialFreq + 1500 - 100
	freqTxHigh := dialFreq + 1500 + 100
	freqTxWindow := freqTxHigh - freqTxLow

	// divide the frequency window into 5 bands
	bandSizeHz := freqTxWindow / 5

	freqBandLow := (frequencyBand(band) - 1) * bandSizeHz
	freqBandHigh := freqBandLow + bandSizeHz

	freqCenter := (freqBandHigh + freqBandLow) / 2

	return uint64(freqTxLow + freqCenter)
}

func transmitWSPRMessage() {
	updateWatchdog()

	log("Transmitting WSPR message...")
	location := wspr.Maidenhead(currentLatitude, currentLongitude)
	log("Callsign:", callsign, "Location:", location, location[:4])
	msg, err := wspr.NewMessage(callsign, location[:4], 37)
	if err != nil {
		log("Error creating WSPR message:", err.Error())
		return
	}

	n, err := msg.WriteSymbols(data[:])
	if err != nil {
		log("error writing WSPR message")
		return
	}

	if err := transmitter.WriteSymbols(data[:n]); err != nil {
		log("error transmitting WSPR message:", err.Error())
		return
	}
}

func transmitTelemetryMessage() {
	updateWatchdog()

	log("Transmitting telemetry message...")
	location := wspr.Maidenhead(currentLatitude, currentLongitude)
	msg, err := u4b.NewMessage(channelCode(), location[:2], int(currentAltitude), int(temperature), int(currentVoltage), int(currentSpeed))
	if err != nil {
		log("Error creating telemetry message:", err.Error())
		return
	}

	n, err := msg.WriteSymbols(data[:])
	if err != nil {
		log("error writing telemetry message")
		return
	}

	if err := transmitter.WriteSymbols(data[:n]); err != nil {
		log("error transmitting telemetry message:", err.Error())
		return
	}
}
