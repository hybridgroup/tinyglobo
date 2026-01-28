package main

var (
	data [256]byte

	// battery voltage in millivolts
	currentVoltage uint32
	// GPS readings
	currentLatitude  float32
	currentLongitude float32
	currentAltitude  int32
	currentSpeed     uint32
)
