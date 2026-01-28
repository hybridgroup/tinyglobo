//go:build tinygo

package main

import (
	"machine"
	"sync"

	"runtime"
	"time"

	"tinygo.org/x/drivers/gps"
)

var (
	ublox      *gps.Device
	currentFix gps.Fix
)

var (
	gpsReset      = machine.GPIO6
	gpsLoadSwitch = machine.GPIO2
)

var (
	gpsStarted     bool
	gpsStopChan    chan struct{}
	gpsFixAcquired bool

	mu            sync.Mutex
	gpsTimeAdjust bool
)

func gpsIsStarted() bool {
	return gpsStarted
}

func gpsHasFix() bool {
	return gpsFixAcquired
}

func gpsTimeAdjusted() bool {
	mu.Lock()
	defer mu.Unlock()

	return gpsTimeAdjust
}

// initialize GPS (called once at startup)
func initGPS() {
	// used to reset GPS module
	gpsReset.Configure(machine.PinConfig{Mode: machine.PinOutput})
	gpsReset.Low()

	// used to control GPS power. high means off. can turn off when not in use.
	gpsLoadSwitch.Configure(machine.PinConfig{Mode: machine.PinOutput})
	gpsLoadSwitch.High()
}

// start GPS reading goroutine
func startGPS() {
	updateWatchdog()
	gpsStarted = true
	gpsFixAcquired = false

	// Power on GPS
	gpsLoadSwitch.Low()
	gpsReset.High()
	time.Sleep(500 * time.Millisecond)

	machine.UART1.Configure(machine.UARTConfig{BaudRate: 9600, RX: machine.UART1_RX_PIN, TX: machine.UART1_TX_PIN})
	u := gps.NewUART(machine.UART1)
	ublox = &u

	gps.SetMessageRatesMinimal(ublox)

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
				log("sentence error:", err)
				continue
			}
		}

		newfix, err := parser.Parse(s)
		if err != nil {
			switch err {
			case gps.ErrUnknownNMEASentence, gps.ErrInvalidNMEASentence, gps.ErrInvalidNMEASentenceLength:
				continue
			default:
				log("parse error:", err)
				continue
			}
		}
		if newfix.Valid {
			currentFix = newfix
			if newfix.Type == gps.GGA || newfix.Type == gps.RMC {
				currentLatitude = newfix.Latitude
				currentLongitude = newfix.Longitude
				currentAltitude = newfix.Altitude
			}

			// adjust time based on GPS time?
			adjustTimeFromGPS(newfix)

			gpsFixAcquired = true
		}
		if !gpsFixAcquired {
			if newfix.Type == gps.GSV {
				log("Satellites in view:", newfix.Satellites)
			}
		}

		time.Sleep(100 * time.Millisecond)
	}
}

// stop GPS reading and power down GPS
func stopGPS() {
	gpsFixAcquired = false

	updateWatchdog()
	close(gpsStopChan)
	time.Sleep(100 * time.Millisecond)

	// Power off GPS
	machine.UART1.Close()
	machine.UART1_TX_PIN.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.UART1_TX_PIN.Low()
	machine.UART1_RX_PIN.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.UART1_RX_PIN.Low()
	gpsReset.Low()
	gpsLoadSwitch.High()

	gpsStarted = false
}

func adjustTimeFromGPS(fix gps.Fix) {
	if !gpsTimeAdjust {
		mu.Lock()
		defer mu.Unlock()

		now := time.Now()
		runtime.AdjustTimeOffset(int64(fix.Time.Sub(now)))
		log("Adjusting system time based on GPS fix from", now.Format("15:04:05"), "to", fix.Time.Format("15:04:05"))

		gpsTimeAdjust = true
	}
}
