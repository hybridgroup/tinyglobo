package main

import (
	"machine"
	"time"

	"github.com/hybridgroup/tinyglobo/powman"
)

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
		time.Sleep(10 * time.Millisecond)

		readBattery()
		println("Battery voltage:", currentVoltage, "mV")
		if currentVoltage > desiredStartingBatteryVoltage {
			break
		}

		println("Battery voltage below desired starting voltage, entering deep sleep...")
		notify(int(StatusIdle))
		time.Sleep(500 * time.Millisecond)
		notify((int(currentVoltage) / 1000) - 1)

		deepSleepForMs(deepSleepDuration)
	}

	startNotification(5 * time.Second)

	initWatchdog()
	initGPS()
	initRadio()
	initSensors()

	for {
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
			// do we still have enough battery to transmit?
			// if not go into deep sleep, the GPS warm start will use less power
			// to obtain a fix next time.
			readBattery()
			if currentVoltage < minTransmitVoltage {
				println("Battery voltage too low for transmission:", currentVoltage, "mV")
				Status = StatusIdle
				stopGPS()

				notify(int(StatusIdle))
				time.Sleep(500 * time.Millisecond)
				notify((int(currentVoltage) / 1000) - 1)

				deepSleepForMs(deepSleepDuration)
			}

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

			// require new GPS fix/time for next transmission, so deep sleep
			// until 4 minutes before the next scheduled transmission
			next := nextScheduledTransmission().Add(-4 * time.Minute)
			deepSleepForMs(uint32(time.Now().Sub(next).Milliseconds()))
		}
	}
}

func deepSleepForMs(ms uint32) {
	machine.LED.Low()
	powman.PinIsolate(uint8(machine.LED))
	if err := powman.SleepForMs(uint64(ms)); err != nil {
		println("Error entering deep sleep:", err.Error())
	}
	time.Sleep(10 * time.Millisecond)
}
