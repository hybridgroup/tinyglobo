//go:build !tinygo

package main

func initBattery() {
}

// calculate in millivolts
func readBattery() uint32 {
	return 32000
}
