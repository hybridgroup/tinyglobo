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
GP0 UART TX     -> GPS RX
GP1 UART RX     -> GPS TX
3V3             -> GPS VCC

### Sensors
GP4 I2C0 SDA    -> HW-290 SDA
GP5 I2C0 SCL    -> HW-290 SCL
3V3             -> HW-290 3V3
GND             -> HW-290 GND

### LAMBDA62 (SX1262)
GP6             -> SX1262 DIO1
GP7             -> SX1262 BUSY (DIO0)
GP8             -> SX1262 TX EN
GP9             -> SX1262 RX EN
GP10 SPI0 CLK   -> SX1262 CLK
GP11 SPI0 CDO   -> SX1262
GP12 SPI0 CDI   -> SX1262
GP13 SPI0 CS    -> SX1262 NSS
3V3             -> SX1262 3V3
GND             -> SX1262 GND


## Samples

### GPS

tinygo flash -target pico ./samples/gps

