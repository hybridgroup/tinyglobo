flash:
	@tinygo flash -size short -target pico2 -ldflags="-X main.callsign=K1ABC -X main.watchdog=enabled" .

build:
	@tinygo build -o tinyglobo.uf2 -size short -target pico2 -ldflags="-X main.callsign=K1ABC -X main.watchdog=enabled" .

verify-gps:
	@tinygo flash -size short -target pico2 -monitor ./verify/gps
