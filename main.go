//go:build tinygo

package main

import (
	"machine"
	"math"
	"time"

	"github.com/hybridgroup/tinyglobo/powman"
)

func main() {
	powman.Init(0)
	machine.InitSerial()

	time.Sleep(3 * time.Second)
	log("*** TinyGlobo 3 starting... ***")

	for {
		initBattery()
		time.Sleep(10 * time.Millisecond)

		readBattery()
		log("Battery voltage:", currentVoltage, "mV")
		if currentVoltage > desiredStartingBatteryVoltage {
			break
		}

		log("Battery voltage below desired starting voltage, entering deep sleep...")
		notify(int(StatusLowBattery))
		time.Sleep(500 * time.Millisecond)
		notify((int(currentVoltage) / 1000) - 1)

		deepSleepForMs(deepSleepDuration)
	}

	//loadStateFromScratch()
	startNotifications(5 * time.Second)

	initWatchdog()
	initSensors()

	for {
		switch {
		case !gpsIsStarted() && !gpsHasFix():
			log("Starting GPS...")
			Status = StatusAcquiringFix
			initGPS()
			go startGPS()
			watchAndWait(5)
			continue

		// wait until we have a fix
		case gpsIsStarted() && !gpsTimeAdjusted():
			readBattery()
			if currentVoltage < criticalLowVoltage {
				log("Battery voltage critical:", currentVoltage, "mV")
				stopGPS()

				stopNotifications()
				notify(int(StatusLowBattery))
				time.Sleep(500 * time.Millisecond)
				notify((int(currentVoltage) / 1000) - 1)

				deepSleepForMs(deepSleepDuration)
			}

			watchAndWait(1)
			continue

		case gpsTimeAdjusted():
			// check if we are in a geofenced area
			if geofenced() {
				log("In geofenced area, delaying transmission.")
				Status = StatusGeofenceBreach
				stopGPS()

				stopNotifications()
				notify(int(StatusGeofenceBreach))
				time.Sleep(500 * time.Millisecond)
				notify(int(StatusGeofenceBreach))

				// wait a half hour to see if we move out of geofenced area
				deepSleepForMs(30 * 60 * 1000)
			}

			// do we have enough battery to transmit?
			// if not go into deep sleep, the GPS warm start will use less power
			// to obtain a fix next time.
			readBattery()
			if currentVoltage < minTransmitVoltage {
				log("Battery voltage too low for transmission:", currentVoltage, "mV")
				Status = StatusIdle
				stopGPS()

				stopNotifications()
				notify(int(StatusLowBattery))
				time.Sleep(500 * time.Millisecond)
				notify((int(currentVoltage) / 1000) - 1)

				deepSleepForMs(deepSleepDuration)
			}

			log("Preparing to transmit...")
			Status = StatusReadyToTransmit
			stopGPS()

			now := time.Now()
			transmit := nextScheduledTransmissionFrom(now)
			if transmit.Sub(now) > 120*time.Second {
				// too early, go back to sleep until 120 seconds before next transmission
				next := transmit.Add(-120 * time.Second)
				ms := next.Sub(now).Milliseconds()
				if ms > 1000 {
					log("Too soon before transmission, deep sleep for", ms, "ms")

					stopNotifications()
					notify(int(StatusTooSoonToTransmit))
					time.Sleep(500 * time.Millisecond)
					notify(int(StatusTooSoonToTransmit))

					deepSleepForMs(uint32(ms - 1000))
				}
			}

			log("Waiting for warmup...")
			waitUntil(transmit.Add(-30 * time.Second))
			initRadio()
			startRadio()
			readSensors()

			log("Waiting for transmission window...")
			waitUntil(transmit)

			Status = StatusTransmitting
			transmitWSPRMessage()

			readBattery()
			log("Battery voltage after WSPR:", currentVoltage, "mV")

			// send the telemetry message 2 minutes after the WSPR message
			waitUntil(transmit.Add(120 * time.Second))
			transmitTelemetryMessage()

			stopRadio()

			Status = StatusIdle
			log("Transmission complete.")

			// require new GPS fix/time for next transmission, so deep sleep
			// until 3 minutes before the next scheduled transmission
			now = time.Now()
			next := nextScheduledTransmissionFrom(now).Add(-3 * time.Minute)
			ms := next.Sub(now).Milliseconds()
			if ms > 0 {
				deepSleepForMs(uint32(ms))
			}
		}
	}
}

func deepSleepForMs(ms uint32) {
	// save gps info into scratch registers
	// saveStateToScratch()

	machine.LED.Low()
	powman.PinIsolate(uint8(machine.LED))
	powman.SleepForMs(uint64(ms))
	time.Sleep(10 * time.Millisecond)
}

func saveStateToScratch() {
	powman.Scratch(0).Set(uint32(Status))
	powman.Scratch(1).Set(math.Float32bits(currentLatitude))
	powman.Scratch(2).Set(math.Float32bits(currentLongitude))
	powman.Scratch(3).Set(uint32(currentAltitude))
	powman.Scratch(4).Set(uint32(currentSpeed))
}

func loadStateFromScratch() {
	Status = StatusType(powman.Scratch(0).Get())
	currentLatitude = math.Float32frombits(powman.Scratch(1).Get())
	currentLongitude = math.Float32frombits(powman.Scratch(2).Get())
	currentAltitude = int32(powman.Scratch(3).Get())
	currentSpeed = uint32(powman.Scratch(4).Get())
}
