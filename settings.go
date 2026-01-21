package main

const (
	minBatteryVoltage             = 2500 // millivolts
	desiredStartingBatteryVoltage = 3600 // millivolts
)

// Things we have to set up via -ldflags:
var (
	// the balloon's callsign, eg "EA1GVK". Obtain from your national amateur radio authority.
	callsign string

	// the channel code, eg "01" or "Q8". obtain from https://traquito.github.io/channelmap/
	code string

	// the channel lane, "1" to "4". obtain from https://traquito.github.io/channelmap/
	lane string

	// minutes past the 10-minute interval to start transmission. obtain from https://traquito.github.io/channelmap/
	offset string
)

// channelCode returns the channel code, eg "01" or "Q8".
func channelCode() string {
	return code
}

// channelLane returns the channel lane as an integer.
func channelLane() int {
	switch lane {
	case "1":
		return 1
	case "2":
		return 2
	case "3":
		return 3
	case "4":
		return 4
	default:
		// should we error out here?
		return 1
	}
}

// transmissionOffset returns the transmission offset in minutes form each 10-minute interval.
func transmissionOffset() int {
	switch offset {
	case "0":
		return 0
	case "2":
		return 2
	case "4":
		return 4
	case "6":
		return 6
	case "8":
		return 8
	default:
		// should we error out here?
		return 0
	}
}
