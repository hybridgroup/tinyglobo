module github.com/hybridgroup/tinyglobo

go 1.25.3

replace tinygo.org/x/drivers => ../tinygo/drivers

require tinygo.org/x/drivers v0.34.1-0.20260108130541-892265b7332b

require (
	github.com/TheThingsNetwork/go-cayenne-lib v1.1.0
	tinygo.org/x/wireless v0.0.0-20260108104403-d628e11a764e
)

require github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect

replace github.com/TheThingsNetwork/go-cayenne-lib => github.com/tinygo-org/go-cayenne-lib v0.0.0-20230116181903-6c670c96018c
