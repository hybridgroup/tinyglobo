flash:
	@tinygo flash -size short -target pico -ldflags="-X main.appEUI='$(shell cat ./keys/appeui)' -X main.devEUI='$(shell cat ./keys/deveui)' -X main.appKey='$(shell cat ./keys/appkey)'" .
