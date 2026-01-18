package main

import (
	"time"
)

var (
	callsign              string
	transmissionFrequency = 10 * time.Minute

	data             [128 + 64]byte
	lastTransmission time.Time
	currentLatitude  float32
	currentLongitude float32
	currentAltitude  int32
)

func main() {
	time.Sleep(3 * time.Second)
	println("*** TinyGlobo 3 starting... ***")

	startNotification(5 * time.Second)

	initWatchdog()

	initBattery()
	for {
		readBattery()
		if readBattery() > desiredStartingBatteryVoltage {
			break
		}

		// wait 15 seconds before checking again
		watchAndWait(30)
	}

	initGPS()
	initRadio()
	initSensors()

	for {
		// check if we have enough battery voltage to transmit
		readBattery()
		// if voltage < desiredBatteryVoltage {
		// 	Status = StatusIdle
		// 	watchAndWait(15)
		// 	continue
		// }

		switch {
		case !gpsIsStarted() && !gpsHasFix():
			println("Starting GPS...")
			Status = StatusAcquiringFix
			go startGPS()
			watchAndWait(5)
			continue

		// wait until we have a fix
		case gpsIsStarted() && !gpsHasFix():
			// TODO: add timeout and restart GPS if needed
			println("Waiting for GPS fix...")
			watchAndWait(1)
			continue

		case gpsHasFix():
			// check if we are in a geofenced area
			if geofenced() {
				println("In geofenced area, delaying transmission.")
				Status = StatusIdle
				stopGPS()

				// wait a half hour to see if we move out of geofenced area
				watchAndWait(1800)
				continue
			}

			println("Preparing to transmit...")
			Status = StatusReadyToTransmit
			stopGPS()

			transmit := nextScheduledTransmission()
			println("Waiting for warmup...")

			waitUntil(transmit.Add(-1 * time.Minute))
			startRadio()
			readSensors()

			println("Waiting for transmission window...")
			waitUntil(transmit)

			Status = StatusTransmitting
			transmitWSPRMessage()
			stopRadio()

			Status = StatusIdle
			println("Transmission complete.")

			// require new GPS fix/time for next transmission
			waitUntil(nextScheduledTransmission().Add(-4 * time.Minute))
		}
	}
}

func failure(err error) {
	for {
		println("FATAL:", err)
		time.Sleep(time.Second)
	}
}
