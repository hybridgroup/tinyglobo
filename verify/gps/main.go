package main

import (
	"machine"
	"time"

	"tinygo.org/x/drivers/gps"
)

func main() {

	reset := machine.GPIO6
	reset.Configure(machine.PinConfig{Mode: machine.PinOutput})
	gpsPowerOn := machine.GPIO3
	gpsPowerOn.Configure(machine.PinConfig{Mode: machine.PinOutput})
	gpsLoadSwitch := machine.GPIO2
	gpsLoadSwitch.Configure(machine.PinConfig{Mode: machine.PinOutput})

	gpsLoadSwitch.High()
	gpsPowerOn.Low()
	reset.Low()

	time.Sleep(5 * time.Second)

	machine.UART1.Configure(machine.UARTConfig{BaudRate: 9600, RX: machine.UART1_RX_PIN, TX: machine.UART1_TX_PIN})

	// Power on GPS
	gpsLoadSwitch.Low()
	reset.High()
	gpsPowerOn.High()
	time.Sleep(500 * time.Millisecond)

	println("GPS UART Example")
	ublox := gps.NewUART(machine.UART1)
	if err := gps.SetMessageRatesMinimal(&ublox); err != nil {
		println("Error setting minimal message rates:", err.Error())
	}

	parser := gps.NewParser()
	var fix gps.Fix
	satelliteFix := false
	for {
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

		fix, err = parser.Parse(s)
		if err != nil {
			switch err {
			case gps.ErrUnknownNMEASentence, gps.ErrInvalidNMEASentence, gps.ErrInvalidNMEASentenceLength:
				continue
			default:
				println("parse error:", err)
				continue
			}
		}
		if fix.Valid {
			satelliteFix = true
			print(fix.Time.Format("15:04:05"))
			print(", lat=")
			print(fix.Latitude)
			print(", long=")
			print(fix.Longitude)
			print(", altitude=", fix.Altitude)
			print(", satellites=", fix.Satellites)
			if fix.Speed != 0 {
				print(", speed=")
				print(fix.Speed)
			}
			if fix.Heading != 0 {
				print(", heading=")
				print(fix.Heading)
			}
			println()
		} else {
			if !satelliteFix {
				if fix.Satellites == 0 {
					println("Waiting for fix...")
				} else {
					println("Waiting for fix...", fix.Satellites, "satellites visible")
				}
			}
		}
		time.Sleep(200 * time.Millisecond)
	}
}
