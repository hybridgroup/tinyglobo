flash:
	@tinygo flash -size short -target pico -ldflags="-X main.callsign=''" .

verify-gps:
	@tinygo flash -size short -target pico -monitor ./verify/gps

