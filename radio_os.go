//go:build !tinygo

package main

import (
	"tinygo.org/x/wireless/fsk4"
)

var (
	radio        NoRadio
	transmitter  *fsk4.FSK4
	radioStarted bool
)

func initRadio() error {
	return nil
}

func startRadio() error {
	return nil
}

func stopRadio() error {
	radioStarted = false

	return nil
}

type NoRadio struct {
}

func (r *NoRadio) Transmit(freq uint64) error {
	return nil
}

func (r *NoRadio) Standby() error {
	return nil
}

func (r *NoRadio) Close() error {
	return nil
}
