flight:
	@tinygo flash -size short -target pico2 -ldflags="-X main.callsign=EA1JZH -X main.code=Q0 -X main.lane=1 -X main.offset=6 -X main.watchdog=enabled" -scheduler=tasks .

flash:
	@tinygo flash -monitor -size short -target pico2 -ldflags="-X main.callsign=EA1JZH -X main.code=Q0 -X main.lane=1 -X main.offset=6 -X main.watchdog=enabled -X main.logging=enabled" -scheduler=tasks .

build:
	@tinygo build -o tinyglobo.uf2 -size short -target pico2 -ldflags="-X main.callsign=EA1JZH -X main.code=Q0 -X main.lane=1 -X main.offset=6 -X main.watchdog=enabled -X main.logging=enabled" -scheduler=tasks .

verify-gps:
	@tinygo flash -size short -target pico2 -monitor ./verify/gps

verify-powman:
	@tinygo flash -size short -target pico2 -monitor -scheduler=tasks ./verify/powman
