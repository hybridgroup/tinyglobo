# TinyGlobo

A small balloon floats into the great big world, towing a RP2040 programmed with TinyGo. This intrepid explorer reports data using LoRaWAN long-range radio.

## Flashing

```
make flash
```

Put your keys into the `/keys` directory as explained there.

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

## Inflation of the balloons

NOTE: keep at least 1 clip attached to each balloon at all times!!!!

You will need 2 - 36 inch transparent balloons.

Use 2 big clips to hold them shut while filling to needed buoyancy.

To calculate:

Full payload of the TinyGlobo / 2 + 5g == Target lift per balloon

2 times clip weight minus target lift per balloon == the weight to be shown on scale when filling balloons.

Use 2 small tie wraps to seal each balloon with a twist.

Connect both balloons together using 2 more small tie wraps.

Connect balloons to payload using fishing line. Make sure the knots are snug and correct.

FLY!!
