package main

import (
	"time"
)

var (
	transmissionOffset = 2 // minutes past the 10-minute interval to start transmission.
)

func waitUntil(when time.Time) {
	sleepDuration := time.Until(when)

	if sleepDuration > 0 {
		// Round down to the nearest second
		sleepDuration -= time.Duration(sleepDuration.Nanoseconds() % 1_000_000_000)
		println("Waiting for", int(sleepDuration.Seconds()), "seconds")
		watchAndWait(int(sleepDuration.Seconds()))
	} else {
		println("Time has already started or passed, proceeding immediately.")
	}
}

func nextScheduledTransmission() time.Time {
	now := time.Now().UTC()

	currentMinute := now.Minute()
	base := (currentMinute / 10) * 10
	candidateMinute := base + transmissionOffset

	// If candidate is in the past or less than 1 minute away, go to next interval
	if candidateMinute <= currentMinute-1 {
		base += 10
		candidateMinute = base + transmissionOffset
	}

	hour := now.Hour()
	day := now.Day()

	// Handle hour rollover
	if candidateMinute >= 60 {
		candidateMinute -= 60
		hour += 1
		if hour >= 24 {
			hour = 0
			day += 1
		}
	}

	return time.Date(now.Year(), now.Month(), day, hour, candidateMinute, 0, 0, time.UTC)
}
