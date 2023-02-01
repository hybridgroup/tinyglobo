package main

import (
	"machine"

	"time"

	"tinygo.org/x/drivers/gps"
)

var fix gps.Fix

func startGPS() {
	machine.UART0.Configure(machine.UARTConfig{BaudRate: 9600})
	ublox := gps.NewUART(machine.UART0)
	parser := gps.NewParser()
	for {
		s, err := ublox.NextSentence()
		if err != nil {
			continue
		}

		newfix, err := parser.Parse(s)
		if err != nil {
			continue
		}
		if newfix.Valid {
			fix = newfix
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
		}
		time.Sleep(200 * time.Millisecond)
	}
}
