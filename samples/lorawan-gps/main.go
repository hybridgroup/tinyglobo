package main

import (
	"machine"

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
	LORAWAN_UPLINK_DELAY_SEC    = 60
)

var (
	spi                        = machine.SPI1
	nssPin, busyPin, dio1Pin   = machine.GP13, machine.GP6, machine.GP7
	rxPin, txLowPin, txHighPin = machine.GP9, machine.GP8, machine.GP8
)

var (
	radio   lora.Radio
	session *lorawan.Session
	otaa    *lorawan.Otaa
)

var (
	encoder cayennelpp.Encoder
)

func loraConnect() error {
	start := time.Now()
	var err error
	for time.Since(start) < LORAWAN_JOIN_TIMEOUT_SEC*time.Second {
		println("Trying to join network")
		err = lorawan.Join(otaa, session)
		if err == nil {
			println("Connected to network !")
			return nil
		}
		println("Join error:", err, "retrying in", LORAWAN_RECONNECT_DELAY_SEC, "sec")
		time.Sleep(time.Second * LORAWAN_RECONNECT_DELAY_SEC)
	}

	err = errors.New("Unable to join Lorawan network")
	println(err.Error())
	return err
}

func failMessage(err error) {
	for {
		println("FATAL:", err)
		time.Sleep(time.Second)
	}
}

func main() {
	time.Sleep(5 * time.Second)
	println("*** Lorawan GPS demo ***")

	spi.Configure(machine.SPIConfig{
		Mode:      0,
		Frequency: 8 * 1e6,
		SDO:       machine.SPI1_SDO_PIN,
		SDI:       machine.SPI1_SDI_PIN,
		SCK:       machine.SPI1_SCK_PIN,
	})

	// Board specific Lorawan initialization
	var err error
	radio, err = SetupLora()
	if err != nil {
		failMessage(err)
	}

	// Required for LoraWan operations
	session = &lorawan.Session{}
	otaa = &lorawan.Otaa{}

	// Connect the lorawan with the Lora Radio device.
	lorawan.UseRadio(radio)

	// Configure AppEUI, DevEUI, APPKey
	setLorawanKeys()

	// Try to connect Lorawan network
	if err := loraConnect(); err != nil {
		failMessage(err)
	}

	go startGPS()

	encoder = cayennelpp.NewEncoder()

	// Try to send an uplink message
	for {
		println("Sleeping for", LORAWAN_UPLINK_DELAY_SEC, "sec")
		time.Sleep(time.Second * LORAWAN_UPLINK_DELAY_SEC)

		payload, err := createPayload()
		if err != nil {
			continue
		}

		if err := lorawan.SendUplink(payload, session); err != nil {
			println("Uplink error:", err)
		} else {
			println("Uplink success, msglen=", len(payload))
			//println(hex.EncodeToString(payload))
		}
	}
}

func createPayload() ([]byte, error) {
	encoder.Reset()
	encoder.AddPresence(1, 1)
	if fix.Valid {
		println(float64(fix.Latitude), float64(fix.Longitude), float64(fix.Altitude))
		encoder.AddGPS(2, float64(fix.Latitude), float64(fix.Longitude), float64(fix.Altitude))
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
