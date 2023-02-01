package main

import (
	"encoding/hex"
	"time"

	"tinygo.org/x/drivers/lora"
	"tinygo.org/x/drivers/lora/lorawan"

	cayennelpp "github.com/TheThingsNetwork/go-cayenne-lib"
)

const (
	joinTimeoutSeconds    = 180
	reconnectDelaySeconds = 15
	uplinkDelaySeconds    = 30
)

var (
	radio   lora.Radio

	session *lorawan.Session = &lorawan.Session{}
	otaa    *lorawan.Otaa = &lorawan.Otaa{}

	encoder cayennelpp.Encoder = cayennelpp.NewEncoder()
)

func lorawanJoin() error {
	// Configure AppEUI, DevEUI, APPKey
	if err := setLorawanKeys(); err != nil {
		return err
	}

	start := time.Now()
	for time.Since(start) < joinTimeoutSeconds*time.Second {
		println("Trying to join LoRaWAN network")
		err := lorawan.Join(otaa, session)
		if err == nil {
			println("Connected to LoRaWAN network!")

			return nil
		}
		println("Join error:", err, "retrying in", reconnectDelaySeconds, "sec")
		time.Sleep(time.Second * reconnectDelaySeconds)
	}

	err := errUnableToJoin
	println(err.Error())
	return err
}

func createPayload() ([]byte, error) {
	encoder.Reset()
	
	// battery voltage
	encoder.AddAnalogInput(1, float64(voltage)/1000)
	
	// BMP180 barometer
	encoder.AddBarometricPressure(2, float64(pressure)/100000)

	// BMP180 temperature
	encoder.AddTemperature(3, float64(temperature)/1000)

	// BMP180 altitude
	encoder.AddAnalogInput(4, float64(altitude))

	// MPU-6050 gyrometer
	encoder.AddAccelerometer(5, float64(ax)/1000000, float64(ay)/1000000, float64(az)/1000000)

	// MPU-6050 gyrometer
	encoder.AddGyrometer(6, float64(ax)/1000000, float64(ay)/1000000, float64(az)/1000000)

	// GPS
	if fix.Valid {
		println(float64(fix.Latitude), float64(fix.Longitude), float64(fix.Altitude))
		encoder.AddGPS(7, float64(fix.Latitude), float64(fix.Longitude), float64(fix.Altitude))
	}

	payload := encoder.Bytes()
	println(hex.EncodeToString(payload))
	return payload, nil
}

var (
	appEUI string
	devEUI string
	appKey string
)

func setLorawanKeys() error {
	if appEUI == "" || devEUI == "" || appKey == "" {
		return errNoKeys
	}

	appEUIData, err := hex.DecodeString(appEUI)
	if err != nil {
		return err
	}
	otaa.SetAppEUI(appEUIData)

	devEUIData, err := hex.DecodeString(devEUI)
	if err != nil {
		return err
	}
	otaa.SetDevEUI(devEUIData)

	appKeyData, err := hex.DecodeString(appKey)
	if err != nil {
		return err
	}
	otaa.SetAppKey(appKeyData)

	lorawan.SetPublicNetwork(true)

	return nil
}
