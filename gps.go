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
	gpsStarted     bool
	gpsStopChan    chan struct{}
	gpsFixAcquired bool
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
	// machine.Watchdog.Update()
	gpsStarted = true

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
			if !gpsFixAcquired {
				now := time.Now()
				println("Adjusting system time based on GPS fix from", now.Format("15:04:05"), "to", newfix.Time.Format("15:04:05"))
				runtime.AdjustTimeOffset(int64(newfix.Time.Sub(now)))
				gpsFixAcquired = true
			}
		}

		time.Sleep(200 * time.Millisecond)
	}
}

// stop GPS reading and power down GPS
func stopGPS() {
	// machine.Watchdog.Update()
	close(gpsStopChan)
	time.Sleep(100 * time.Millisecond)

	// Power off GPS
	gpsReset.Low()
	gpsLoadSwitch.High()

	gpsStarted = false
	gpsFixAcquired = false
}
