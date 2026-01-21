package main

import (
	"time"
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

func nextScheduledTransmissionAt(now time.Time) time.Time {
	currentMinute := now.Minute()
	base := (currentMinute / 10) * 10
	candidateMinute := base + transmissionOffset()

	hour := now.Hour()
	day := now.Day()

	// If candidate time is not in the future, move to next interval
	if candidateMinute <= currentMinute {
		base += 10
		candidateMinute = base + transmissionOffset()
	}

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

// Keep the original for production use
func nextScheduledTransmission() time.Time {
	return nextScheduledTransmissionAt(time.Now().UTC())
}
