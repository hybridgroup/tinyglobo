flash:
	@tinygo flash -size short -target pico -ldflags="-X main.callsign='K1ABC'" .

build:
	@tinygo build -o tinyglobo.uf2 -size short -target pico -ldflags="-X main.callsign='K1ABC'" .

verify-gps:
	@tinygo flash -size short -target pico -monitor ./verify/gps

