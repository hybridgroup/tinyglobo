flash:
	@tinygo flash -size short -target pico2 -ldflags="-X main.callsign=K1ABC -X main.code=Q9 -X main.lane=3 -X main.offset=2 -X main.watchdog=enabled" -scheduler=tasks .

build:
	@tinygo build -o tinyglobo.uf2 -size short -target pico2 -ldflags="-X main.callsign=K1ABC -X main.code=Q9 -X main.lane=3 -X main.offset=2 -X main.watchdog=enabled" -scheduler=tasks .

verify-gps:
	@tinygo flash -size short -target pico2 -monitor ./verify/gps
