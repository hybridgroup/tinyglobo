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

	// configure watchdog for 3 second timeout
	// machine.Watchdog.Configure(machine.WatchdogConfig{
	// 	TimeoutMillis: 3000,
	// })
	// machine.Watchdog.Start()

	machine.I2C0.Configure(machine.I2CConfig{})
	machine.InitADC()

	initGPS()
	initRadio()
	initBattery()
	initSensors()

	for {
		readBattery()

		// TODO: check if we have enough battery voltage to transmit

		switch {
		// ensure at least 10 minutes between transmissions
		case time.Since(lastTransmission) < transmissionFrequency:
			println("Last transmission was too recent, waiting...")
			watchAndWait(15)
			continue

		case !gpsStarted:
			println("Starting GPS...")
			// machine.Watchdog.Update()
			go startGPS()
			watchAndWait(15)
			continue

		// wait until we have a fix
		case !currentFix.Valid:
			println("Waiting for GPS fix...")
			watchAndWait(15)
			continue

		// only transmit on even numbered minutes at exactly 5 second mark
		case time.Now().Minute()%2 != 0:
			println("Preparing to transmit...")

			stopGPS()
			startRadio()
			readSensors()

			// sleep until the next minute that is even numbered
			now := time.Now()
			extra := now.Minute() % 2
			sleepDuration := time.Until(time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute()+extra, 0, 0, now.Location()))
			if sleepDuration > 0 {
				// round down to nearest second
				sleepDuration -= time.Duration(sleepDuration.Nanoseconds() % 1_000_000_000)
				println("Waiting until transmission window...")
				watchAndWait(int(sleepDuration.Seconds()))
			}
			transmitWSPRMessage()
			stopRadio()

			// require new GPS fix for next transmission
			currentFix.Valid = false // require new fix for next transmission
		}
	}
}

func transmitWSPRMessage() {
	lastTransmission = time.Now()
	// machine.Watchdog.Update()

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

func watchAndWait(seconds int) {
	for i := 0; i < seconds; i++ {
		// machine.Watchdog.Update()
		time.Sleep(time.Second)
	}
}
