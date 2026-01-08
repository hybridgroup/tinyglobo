# TinyGlobo

A pico balloon floats into the great big world, towing a [Raspberry Pi RP2040 Pico](https://www.raspberrypi.com/products/raspberry-pi-pico/) programmed with [TinyGo](https://tinygo.org). 

This intrepid explorer reports data using 20 meter band radio using the WSPR protocol.

TinyGlobo 3 dashboard coming soon...

## Previous flights

Check out the [flight](./flights.md) page for info about TinyGlobo 1-2

## Flashing

```
make flash
```

## Callsign

More info here...

## Hardware

![TinyGlobo 3 board](./images/tinyglobo-3-board.jpeg)

TinyGlobo 3 is uses the [Traquito Jetpack Tracker](https://traquito.github.io/tracker/) which consists of:

- RP2040 Pico
- Si5351 radio
- GPS
- Solar power system

## Wiring

### GPS

| RP2040 Pin | GPS Pin |
|------------|---------|
| GPIO8 UART1 TX | RX |
| GPIO9 UART1 RX | TX |
| GPIO3 | VCC |
| GPIO6 | RESET |
| GPIO2 | LOAD (pull low to use) |

### Si5351

| RP2040 Pin | Si5351 Pin |
|------------|---------|
| GPIO4 I2C0 SDA | SDA |
| GPIO5 I2C0 SCL | SCL |
| GPIO28 | Load (pull low to use) |

### Power system

| RP2040 Pin | HW-290 Pin |
|------------|---------|
| GP4 I2C0 SDA | SDA |
| GP5 I2C0 SCL | SCL |
| 3V3 | 3V3 |
| GND | GND |

## Samples

### GPS

tinygo flash -target pico ./verify/gps

