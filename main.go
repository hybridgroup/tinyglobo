package main

import (
	"machine"
	"time"
)

var (
	callsign              string
	transmissionFrequency = 10 * time.Minute

	data             [128 + 64]byte
	lastTransmission time.Time
)

func main() {
	time.Sleep(3 * time.Second)
	println("*** TinyGlobo 3 starting... ***")

	startNotification(5 * time.Second)

	initWatchdog()

	machine.InitADC()
	initBattery()
	for {
		readBattery()
		if voltage > desiredBatteryVoltage {
			break
		}
	}

	machine.I2C0.Configure(machine.I2CConfig{})

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
		case !gpsStarted && !gpsFixAcquired:
			println("Starting GPS...")
			Status = StatusAcquiringFix
			go startGPS()
			watchAndWait(5)
			continue

		// wait until we have a fix
		case gpsStarted && !gpsFixAcquired:
			// TODO: add timeout and restart GPS if needed
			println("Waiting for GPS fix...")
			watchAndWait(1)
			continue

		case gpsFixAcquired:
			// check if we are in a geofenced area
			if geofenced() {
				println("In geofenced area, delaying transmission.")
				Status = StatusIdle
				stopGPS()

				gpsFixAcquired = false
				gpsTimeAdjusted = false

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
			gpsFixAcquired = false
			gpsTimeAdjusted = false
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
