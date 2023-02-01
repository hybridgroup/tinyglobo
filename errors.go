package main

import "errors"

var (
	errUnableToJoin = errors.New("Unable to join LoRaWAN network")
	errNoKeys = errors.New("No LoRaWAN keys provided")
)