package main

var logging string

func loggingEnabled() bool {
	if logging == "" {
		return false
	}

	return true
}

func log(anything ...interface{}) {
	if loggingEnabled() {
		for _, v := range anything {
			print(v, " ")
		}
		println()
	}
}
