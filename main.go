package main

import (
	"machine"
	"time"

	"tinygo.org/x/wireless/wspr"
)

var (
	callsign string

	data             [162]byte
	lastTransmission time.Time
)

func main() {
	machine.UART1.Configure(machine.UARTConfig{BaudRate: 9600, RX: machine.UART1_RX_PIN, TX: machine.UART1_TX_PIN})
	machine.I2C0.Configure(machine.I2CConfig{})

	time.Sleep(3 * time.Second)
	println("*** TinyGlobo 3 starting... ***")

	// configure watchdog for 30 second timeout
	machine.Watchdog.Configure(machine.WatchdogConfig{
		TimeoutMillis: 30000,
	})

	initGPS()
	go startGPS()

	if err := initRadio(); err != nil {
		failure(err)
	}
	startRadio()

	frequency := transmitter.GetBaseFrequency()
	println("Transmitting on frequency", frequency, "Hz")

	startBattery()
	startSensors()

	for {
		readBattery()

		// wait until we have a fix
		if !currentFix.Valid {
			println("Waiting for GPS fix...")
			machine.Watchdog.Update()
			time.Sleep(15 * time.Second)
			continue
		}

		if time.Since(lastTransmission) < time.Minute*5 {
			println("Last transmission was too recent, waiting...")
			machine.Watchdog.Update()
			time.Sleep(15 * time.Second)
			continue
		}

		// only transmit on even numbered minutes at exactly 5 second mark
		if currentFix.Time.Minute()%2 != 0 {
			println("Waiting for even numbered minute...")
			continue
		}
		sleepDuration := time.Until(time.Date(currentFix.Time.Year(), currentFix.Time.Month(), currentFix.Time.Day(), currentFix.Time.Hour(), currentFix.Time.Minute(), 5, 0, currentFix.Time.Location()))
		if sleepDuration > 0 {
			println("Waiting until 5 seconds past the minute...")
			time.Sleep(sleepDuration)
		}

		readSensors()
		transmitWSPRMessage()
	}
}

func transmitWSPRMessage() {
	lastTransmission = time.Now()

	println("Transmitting WSPR message...")
	msg, err := wspr.NewMessage(callsign, wspr.Maidenhead(float64(currentFix.Latitude), float64(currentFix.Longitude)), 21)
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
