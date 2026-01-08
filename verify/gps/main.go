package main

import (
	"machine"
	"time"

	"tinygo.org/x/drivers/gps"
)

func main() {
	machine.UART1.Configure(machine.UARTConfig{BaudRate: 9600, RX: machine.UART1_RX_PIN, TX: machine.UART1_TX_PIN})
	time.Sleep(time.Second)

	reset := machine.GPIO6
	reset.Configure(machine.PinConfig{Mode: machine.PinOutput})
	gpsPowerOn := machine.GPIO3
	gpsPowerOn.Configure(machine.PinConfig{Mode: machine.PinOutput})
	gpsLoadSwitch := machine.GPIO2
	gpsLoadSwitch.Configure(machine.PinConfig{Mode: machine.PinOutput})

	// Power on GPS
	gpsLoadSwitch.Low()
	reset.High()
	gpsPowerOn.High()
	time.Sleep(500 * time.Millisecond)

	time.Sleep(2 * time.Second) // wait for console to start
	println("GPS UART Example")
	ublox := gps.NewUART(machine.UART1)
	parser := gps.NewParser()
	var fix gps.Fix
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
			println("Waiting for fix...")
		}
		time.Sleep(200 * time.Millisecond)
	}
}
