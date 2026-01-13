package main

import (
	"machine"
	"time"

	"tinygo.org/x/drivers/si5351"
	"tinygo.org/x/wireless/fsk4"
)

var (
	radio        Si5351Radio
	transmitter  *fsk4.FSK4
	radioStarted bool
)

var (
	radioLoadSwitch = machine.GPIO28
)

func initRadio() error {
	radioLoadSwitch.Configure(machine.PinConfig{Mode: machine.PinOutput})
	radioLoadSwitch.High()
	time.Sleep(100 * time.Millisecond)

	return nil
}

func startRadio() error {
	updateWatchdog()
	radioLoadSwitch.Low()
	time.Sleep(100 * time.Millisecond)

	dev := si5351.New(machine.I2C0)
	cnf := si5351.Config{
		Capacitance:   si5351.CrystalLoad10PF,
		CrystalOutput: 26_000_000,
	}
	if err := dev.Configure(cnf); err != nil {
		return err
	}

	updateWatchdog()
	dev.SetFrequency(si5351.Clock0, 14_097_060)
	dev.SetDriveStrength(si5351.Clock0, si5351.DriveStrength8MA)
	dev.SetClockPower(si5351.Clock0, true)
	dev.EnableOutput(si5351.Clock0, true)

	// Fan out and invert the first clock signal for a
	// 180-degree phase shift on second clock
	dev.SetClockFanout(si5351.FanoutMultisynth, true)
	dev.SetClockSource(si5351.Clock1, si5351.ClockSourceMS0)
	dev.SetClockInvert(si5351.Clock1, true)

	dev.SetDriveStrength(si5351.Clock1, si5351.DriveStrength8MA)
	dev.SetClockPower(si5351.Clock1, true)
	dev.EnableOutput(si5351.Clock1, true)

	updateWatchdog()
	time.Sleep(500 * time.Millisecond)

	dev.EnableOutput(si5351.Clock0, false)
	dev.EnableOutput(si5351.Clock1, false)

	radio.device = dev

	transmitter = fsk4.NewFSK4(&radio, 14_097_060, 146, 682)
	transmitter.Configure()

	radioStarted = true

	return nil
}

func stopRadio() error {
	updateWatchdog()
	radioLoadSwitch.High()
	time.Sleep(100 * time.Millisecond)

	radioStarted = false

	return nil
}

type Si5351Radio struct {
	device *si5351.Device
}

func (r *Si5351Radio) Transmit(freq uint64) error {
	updateWatchdog()

	if err := r.device.SetRawFrequency(si5351.Clock0, si5351.Frequency(freq)); err != nil {
		return err
	}

	r.device.EnableOutput(si5351.Clock0, true)
	r.device.EnableOutput(si5351.Clock1, true)

	return nil
}

func (r *Si5351Radio) Standby() error {
	updateWatchdog()

	r.device.EnableOutput(si5351.Clock0, false)
	r.device.EnableOutput(si5351.Clock1, false)

	return nil
}

func (r *Si5351Radio) Close() error {
	return nil
}
