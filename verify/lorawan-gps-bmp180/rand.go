// NOT true random number generation on RP2040 but good enough for quick demos.
// Do not use in important production systems.
// Seriously, don't.
package main

import (
	"machine"

	"crypto/rand"
)

func init() {
	rand.Reader = &reader{}
}

type reader struct {}

func (r *reader) Read(b []byte) (n int, err error) {
	if len(b) == 0 {
		return
	}

	var randomByte uint32
	for i := range b {
		if i%4 == 0 {
			randomByte, err = machine.GetRNG()
			if err != nil {
				return n, err
			}
		} else {
			randomByte >>= 8
		}
		b[i] = byte(randomByte)
	}

	return len(b), nil
}
