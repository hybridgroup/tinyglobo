package main

import (
	"machine"
	"time"

	"github.com/hybridgroup/tinyglobo/powman"
)

// duration to sleep in deep sleep mode while waiting for battery to charge (milliseconds)
const deepSleepDuration = 30_000

var (
	data [256]byte

	currentLatitude  float32
	currentLongitude float32
	currentAltitude  int32
	// battery voltage in millivolts
	currentVoltage uint32
	currentSpeed   uint32
)

func main() {
	powman.Init(0)
	machine.InitSerial()

	time.Sleep(3 * time.Second)
	println("*** TinyGlobo 3 starting... ***")

	for {
		initBattery()
		val := readBattery()
		println("Battery voltage:", val, "mV")
		if val > desiredStartingBatteryVoltage {
			break
		}

		machine.LED.Low()
		powman.PinIsolate(uint8(machine.LED))
		if err := powman.SleepForMs(deepSleepDuration); err != nil {
			println("Error entering deep sleep:", err.Error())
		}
		time.Sleep(10 * time.Millisecond)
	}

	startNotification(5 * time.Second)

	initWatchdog()
	initGPS()
	initRadio()
	initSensors()

	for {
		// TODO: check if we have enough battery voltage to transmit
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

			readBattery()
			println("Battery voltage after WSPR:", currentVoltage, "mV")

			// send the telemetry message 2 minutes after the WSPR message
			waitUntil(transmit.Add(2 * time.Minute))
			transmitTelemetryMessage()

			stopRadio()

			Status = StatusIdle
			println("Transmission complete.")

			// require new GPS fix/time for next transmission
			waitUntil(nextScheduledTransmission().Add(-4 * time.Minute))
		}
	}
}
