package main

import (
	"encoding/hex"
	"errors"
	"time"

	"tinygo.org/x/drivers/lora"
	"tinygo.org/x/drivers/lora/lorawan"

	cayennelpp "github.com/TheThingsNetwork/go-cayenne-lib"
)

const (
	LORAWAN_JOIN_TIMEOUT_SEC    = 180
	LORAWAN_RECONNECT_DELAY_SEC = 15
	LORAWAN_UPLINK_DELAY_SEC    = 30
)

var (
	radio   lora.Radio

	session *lorawan.Session = &lorawan.Session{}
	otaa    *lorawan.Otaa = &lorawan.Otaa{}
	encoder cayennelpp.Encoder = cayennelpp.NewEncoder()
)

func loraConnect() error {
	// Configure AppEUI, DevEUI, APPKey
	setLorawanKeys()

	start := time.Now()
	for time.Since(start) < LORAWAN_JOIN_TIMEOUT_SEC*time.Second {
		println("Trying to join network")
		err := lorawan.Join(otaa, session)
		if err == nil {
			println("Connected to network!")

			return nil
		}
		println("Join error:", err, "retrying in", LORAWAN_RECONNECT_DELAY_SEC, "sec")
		time.Sleep(time.Second * LORAWAN_RECONNECT_DELAY_SEC)
	}

	err := errors.New("Unable to join Lorawan network")
	println(err.Error())
	return err
}

func createPayload() ([]byte, error) {
	encoder.Reset()
	
	encoder.AddAnalogInput(1, float64(voltage)/1000)
	
	encoder.AddBarometricPressure(2, float64(pressure)/100000)
	encoder.AddTemperature(2, float64(temperature)/1000)
	encoder.AddAnalogInput(2, float64(altitude)/1000)

	encoder.AddAccelerometer(3, float64(ax)/1000000, float64(ay)/1000000, float64(az)/1000000)
	encoder.AddGyrometer(3, float64(ax)/1000000, float64(ay)/1000000, float64(az)/1000000)

	if fix.Valid {
		println(float64(fix.Latitude), float64(fix.Longitude), float64(fix.Altitude))
		encoder.AddGPS(4, float64(fix.Latitude), float64(fix.Longitude), float64(fix.Altitude))
	}

	payload := encoder.Bytes()
	println(hex.EncodeToString(payload))
	return payload, nil
}

// These are sample keys, so the example builds
// Either change here, or create a new go file and use customkeys build tag
func setLorawanKeys() {
	otaa.SetAppEUI([]uint8{0x70, 0xB3, 0xD5, 0x7E, 0xD0, 0x04, 0xA9, 0x12})
	otaa.SetDevEUI([]uint8{0x70, 0xB3, 0xD5, 0x7E, 0xD0, 0x05, 0x96, 0xEE})
	otaa.SetAppKey([]uint8{0xD8, 0x98, 0x88, 0xEF, 0xF0, 0xB2, 0x61, 0xCB, 0x4D, 0x57, 0x03, 0xAD, 0xD5, 0x87, 0xEF, 0x2B})
}

// AT+ID=AppEui, 70B3D57ED004A912
// AT+ID=DevEui, 70B3D57ED005933B -> 70B3D57ED00596EE
// AT+KEY=APPKEY, 6757BB981D0E2671F40F534F6E4CD87F -> D89888EFF0B261CB4D5703ADD587EF2B
