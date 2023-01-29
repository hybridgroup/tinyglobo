# TinyGlobo

A small balloon floats into the great big world, towing a RP2040 programmed with TinyGo. This intrepid explorer reports data using LoRaWAN long-range radio.

## Flashing

tinygo flash -target pico .


## Wiring

TinyGlobo 1 consists of:

- RP2040 Pico
- LAMBDA62 LoRa radio
- UBlox 6M GPS
- HW-290 multifunction I2C board
    - MPU6050
    - BMP180
    - HMC5883L

### GPS

| RP2040 Pin | GPS Pin |
|------------|---------|
| GP0 UART TX | RX |
| GP1 UART RX | TX |
| 3V3 | VCC |

### Sensors

| RP2040 Pin | HW-290 Pin |
|------------|---------|
| GP4 I2C0 SDA | SDA |
| GP5 I2C0 SCL | SCL |
| 3V3 | 3V3 |
| GND | GND |

### LAMBDA62 (SX1262)

| RP2040 Pin | LAMBDA62 Pin |
|------------|---------|
| GP6 | DIO1 |
| GP7 | DIO0 |
| GP8 | TX_SWITCH |
| GP9 | RX_SWITCH |
| GP10 | SCLK |
| GP11 SPI0 CDO | SDI |
| GP12 SPI0 CDI | SDO |
| GP13 SPI0 CS | nSEL |
| 3V3 | 3V3 |
| GND | GND |



## Samples

### GPS

tinygo flash -target pico ./samples/gps

