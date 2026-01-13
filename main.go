package main

import (
	"machine"
	"time"

	"tinygo.org/x/wireless/wspr"
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
		// ensure at least 10 minutes between transmissions
		case time.Since(lastTransmission) < transmissionFrequency:
			println("Last transmission was too recent, waiting...")
			watchAndWait(15)
			continue

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
			watchAndWait(5)
			continue

		case gpsFixAcquired:
			println("Preparing to transmit...")
			Status = StatusReadyToTransmit

			stopGPS()
			startRadio()
			readSensors()

			// only transmit on even numbered minutes
			now := time.Now()
			minute := now.Minute()
			nextEvenMinute := (minute + 2 - (minute % 2)) % 60
			nextTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), nextEvenMinute, 0, 0, now.Location())

			// If the next even minute is less than 1 minute away, skip to the following even minute
			if nextTime.Sub(now) < time.Minute {
				nextEvenMinute = (nextEvenMinute + 2) % 60
				nextTime = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), nextEvenMinute, 0, 0, now.Location())
			}

			sleepDuration := time.Until(nextTime)
			if sleepDuration > 0 {
				// round down to nearest second
				sleepDuration -= time.Duration(sleepDuration.Nanoseconds() % 1_000_000_000)
				println("Waiting until transmission window...")
				watchAndWait(int(sleepDuration.Seconds()))
			}

			Status = StatusTransmitting
			transmitWSPRMessage()
			stopRadio()

			Status = StatusIdle
			println("Transmission complete.")

			// require new GPS fix for next transmission
			gpsFixAcquired = false
		}
	}
}

func transmitWSPRMessage() {
	lastTransmission = time.Now()
	updateWatchdog()

	println("Transmitting WSPR message...")
	location := wspr.Maidenhead(float64(currentFix.Latitude), float64(currentFix.Longitude))
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

	println("Transmitting WSPR message with", n, "symbols")
	if err := transmitter.WriteSymbols(data[:n]); err != nil {
		println("error transmitting WSPR message:", err.Error())
		return
	}
}

func failure(err error) {
	for {
		println("FATAL:", err)
		time.Sleep(time.Second)
	}
}
