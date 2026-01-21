package main

import (
	"machine"
	"time"

	"tinygo.org/x/drivers/gps"
)

func main() {

	reset := machine.GPIO6
	reset.Configure(machine.PinConfig{Mode: machine.PinOutput})
	gpsLoadSwitch := machine.GPIO2
	gpsLoadSwitch.Configure(machine.PinConfig{Mode: machine.PinOutput})

	gpsLoadSwitch.High()
	reset.Low()

	time.Sleep(5 * time.Second)
	println("GPS UART Example")

	for i := 0; i < 3; i++ {
		machine.UART1.Configure(machine.UARTConfig{BaudRate: 9600, RX: machine.UART1_RX_PIN, TX: machine.UART1_TX_PIN})

		// Power on GPS
		gpsLoadSwitch.Low()
		reset.High()
		time.Sleep(500 * time.Millisecond)

		ublox := gps.NewUART(machine.UART1)
		if err := gps.SetMessageRatesMinimal(&ublox); err != nil {
			println("Error setting minimal message rates:", err.Error())
		}

		parser := gps.NewParser()
		var fix gps.Fix
		for {
			s, err := ublox.NextSentence()
			println(s)
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
				break
			} else {
				if fix.Type == gps.GSV {
					println("Satellites in view:", fix.Satellites)
				}
			}
			time.Sleep(200 * time.Millisecond)
		}

		println("Shutting down GPS...")

		// Power off GPS
		machine.UART1.Close()
		machine.UART1_TX_PIN.Configure(machine.PinConfig{Mode: machine.PinOutput})
		machine.UART1_TX_PIN.Low()
		machine.UART1_RX_PIN.Configure(machine.PinConfig{Mode: machine.PinOutput})
		machine.UART1_RX_PIN.Low()
		reset.Low()
		gpsLoadSwitch.High()

		time.Sleep(1 * time.Minute)
	}

	for {
		println("GPS test complete.")
		time.Sleep(10 * time.Second)
	}
}
