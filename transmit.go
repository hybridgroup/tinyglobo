package main

import (
	"time"

	"tinygo.org/x/wireless/wspr"
)

const dialFreq = 14_095_600 // in Hz

// 1, 2, 4, 5 using channel map that skips 3 to avoid interference.
// TODO: set this from config or command line argument.
var channelBand int = 1

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
	lastTransmission = time.Now()
	updateWatchdog()

	println("Transmitting WSPR message...")
	location := wspr.Maidenhead(currentFix.Latitude, currentFix.Longitude)
	println("Callsign:", callsign, "Location:", location, location[:4])
	msg, err := wspr.NewMessage(callsign, location[:4], 37)
	if err != nil {
		println("Error creating WSPR message:", err.Error())
		return
	}

	n, err := msg.WriteSymbols(data[:])
	if err != nil {
		println("error writing WSPR message")
		return
	}

	if err := transmitter.WriteSymbols(data[:n]); err != nil {
		println("error transmitting WSPR message:", err.Error())
		return
	}
}
