package main

import (
	"machine"
	"time"

	"tinygo.org/x/drivers/si5351"
	"tinygo.org/x/wireless/fsk4"
)

func initRadio() (*fsk4.FSK4, error) {
	machine.I2C0.Configure(machine.I2CConfig{})
	gpsLoadSwitch := machine.GPIO28
	gpsLoadSwitch.Configure(machine.PinConfig{Mode: machine.PinOutput})
	gpsLoadSwitch.High()
	time.Sleep(100 * time.Millisecond)
	gpsLoadSwitch.Low()
	time.Sleep(100 * time.Millisecond)

	dev := si5351.New(machine.I2C0)
	cnf := si5351.Config{
		Capacitance:   si5351.CrystalLoad10PF,
		CrystalOutput: 26_000_000,
	}
	if err := dev.Configure(cnf); err != nil {
		return nil, err
	}

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

	time.Sleep(500 * time.Millisecond)

	dev.EnableOutput(si5351.Clock0, false)
	dev.EnableOutput(si5351.Clock1, false)

	f := fsk4.NewFSK4(&Si5351Radio{device: dev}, 14_097_060, 146, 682)
	f.Configure()

	return f, nil
}

type Si5351Radio struct {
	device *si5351.Device
}

func (r *Si5351Radio) Transmit(freq uint64) error {
	if err := r.device.SetRawFrequency(si5351.Clock0, si5351.Frequency(freq)); err != nil {
		return err
	}

	r.device.EnableOutput(si5351.Clock0, true)
	r.device.EnableOutput(si5351.Clock1, true)

	return nil
}

func (r *Si5351Radio) Standby() error {
	r.device.EnableOutput(si5351.Clock0, false)
	r.device.EnableOutput(si5351.Clock1, false)

	return nil
}

func (r *Si5351Radio) Close() error {
	return nil
}
