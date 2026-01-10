package main

import (
	"machine"

	"runtime"
	"time"

	"tinygo.org/x/drivers/gps"
)

var (
	ublox      *gps.Device
	currentFix gps.Fix
)

var (
	gpsReset          = machine.GPIO6
	gpsLoadSwitch     = machine.GPIO2
	gpsBatteryPowerOn = machine.GPIO3
)

var (
	gpsStarted         bool
	gpsStopChan        chan struct{}
	lastTimeAdjustment time.Time
)

// initialize GPS (called once at startup)
func initGPS() {
	// used to reset GPS module
	gpsReset.Configure(machine.PinConfig{Mode: machine.PinOutput})
	gpsReset.Low()

	// used to control GPS power. high means off. can turn off when not in use.
	gpsLoadSwitch.Configure(machine.PinConfig{Mode: machine.PinOutput})
	gpsLoadSwitch.High()

	// used to control GPS battery power. leave on for warm starts.
	gpsBatteryPowerOn.Configure(machine.PinConfig{Mode: machine.PinOutput})
	gpsBatteryPowerOn.High()
}

// start GPS reading goroutine
func startGPS() {
	gpsStarted = true

	// Power on GPS
	gpsLoadSwitch.Low()
	gpsReset.High()
	time.Sleep(500 * time.Millisecond)

	u := gps.NewUART(machine.UART1)
	ublox = &u
	parser := gps.NewParser()
	gpsStopChan = make(chan struct{})
	for {
		select {
		case <-gpsStopChan:
			return
		default:
		}

		s, err := ublox.NextSentence()
		if err != nil {
			switch err {
			case gps.ErrUnknownNMEASentence, gps.ErrInvalidNMEASentence, gps.ErrInvalidNMEASentenceLength:
				continue
			default:
				println("sentence error:", err)
				continue
			}
		}

		newfix, err := parser.Parse(s)
		if err != nil {
			switch err {
			case gps.ErrUnknownNMEASentence, gps.ErrInvalidNMEASentence, gps.ErrInvalidNMEASentenceLength:
				continue
			default:
				println("parse error:", err)
				continue
			}
		}
		if newfix.Valid {
			currentFix = newfix

			// adjust time based on GPS time
			if newfix.Time.Sub(lastTimeAdjustment) > time.Minute*10 {
				now := time.Now()
				println("Adjusting system time based on GPS fix from", now.Format("15:04:05"), "to", newfix.Time.Format("15:04:05"))
				runtime.AdjustTimeOffset(int64(newfix.Time.Sub(now)))
				lastTimeAdjustment = newfix.Time
			}

			// print(currentFix.Time.Format("15:04:05"))
			// print(", lat=")
			// print(currentFix.Latitude)
			// print(", long=")
			// print(currentFix.Longitude)
			// print(", altitude=", currentFix.Altitude)
			// print(", satellites=", currentFix.Satellites)
			// if currentFix.Speed != 0 {
			// 	print(", speed=")
			// 	print(currentFix.Speed)
			// }
			// if currentFix.Heading != 0 {
			// 	print(", heading=")
			// 	print(currentFix.Heading)
			// }
			// println()
		}
		time.Sleep(200 * time.Millisecond)
	}
}

// stop GPS reading and power down GPS
func stopGPS() {
	close(gpsStopChan)
	time.Sleep(100 * time.Millisecond)

	// Power off GPS
	gpsLoadSwitch.High()
	gpsReset.Low()

	gpsStarted = false
}
