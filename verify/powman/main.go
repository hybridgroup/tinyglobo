//go:build rp2350

package main

import (
	"machine"
	"time"

	"github.com/hybridgroup/tinyglobo/powman"
)

const sleepDuration = 10_000

func main() {
	// Initialize POWMAN with current time (0ms at boot)
	powman.Init(0)
	println("start program")
	// Main loop: do work, sleep, repeat
	for {
		// Configure LED for work indication
		led := machine.LED
		led.Configure(machine.PinConfig{Mode: machine.PinOutput})
		doWork(led)
		println("work done")
		time.Sleep(5 * time.Second)
		println("LED off")
		led.Low()

		time.Sleep(5 * time.Second)
		println("go to sleep :)")
		// isolate LED for lowest power consumption.
		// If this is not done the deep sleep will cause SIO to light the LED up, defeating low power mode purpose.
		// Keep an eye out for weird peripheral behaviour after deep sleep.
		// Serial is glitchy and may eventually break unless lots of sleep is added.
		// Least glitchy results by increasing serialflush sleep.
		powman.PinIsolate(uint8(machine.LED))
		// Enter low power sleep.
		powman.SleepForMs(sleepDuration)
		time.Sleep(time.Millisecond)
		// Serial will have been completely deinitialized after deep sleep.
		machine.InitSerial()
		time.Sleep(time.Millisecond) // wait a bit for serial to start up.
		println("good morning sunshine :)")
	}
}

// doWork simulates some work by blinking the LED
func doWork(led machine.Pin) {
	// Blink LED 5 times to indicate we're doing work
	for i := 0; i < 5; i++ {
		led.High()
		time.Sleep(100 * time.Millisecond)
		led.Low()
		time.Sleep(100 * time.Millisecond)
	}
}
