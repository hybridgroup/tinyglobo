# Assembly

The payload consists of 3 main subsystems:

- Airframe
- Control module
- Solar module

## Airframe

Space Cluster with 2 connectors and 3 struts cut down to 23cm length. The connectors each need a notch cut at approx. 1.2cm from top to serve as anchor for the solar panels.

## Control module

- [ ] carefully remove power regulator from Traquito.

- [ ] solder pico2 to Traquito board at corners to hold together

- [ ] solder all other pins on pico2/traquito from instructions except for pins `3V3_OUT` and `GP5`

- [ ] solder jumper across pins `3V3_OUT` and `GP5` while soldering to Traquito

- [ ] solder male dupont gnd cable (8cm) to pico2 gnd pin just below Traquito on same side as solar +

- [ ] solder male dupont positive solar cable (8cm) to Traquito

- [ ] solder female dupont negative solar cable to Traquito

- [ ] solder GPS antenna cable (14cm) to Traquito on solar + side

- [ ] attach fishing line horizontal loop on pico2

- [ ] solder top antenna wire on solar ground side

- [ ] solder bottom antenna wire on solar + side

- [ ] tie control module to top strut of airframe

## Solar module

The solar module consists of the 2 solar panels, and a super capacitor to hold a small amount of charge needed for system functions.

### Solar panels

***NOTE** the panels must be oriented to be connected together in parallel. In other words, the + both will be on the top side, and the - for both will be on the bottom.

#### Panel 1

- [ ] melt spot for connect of + and - on solar cell at low temperature. put the + on the left side, and the - on the right side

- [ ] drop small bead of solder on each spot at normal temperature

- [ ] solder + and - wires to each spot

- [ ] anchor each wire with krylon tape

#### Panel 2

- [ ] melt spot for connect of + and - on solar cell at low temperature. put the + on the right side, and the - on the left side

- [ ] drop small bead of solder on each spot at normal temperature

- [ ] solder + and - wires to each spot

- [ ] anchor each wire with krylon tape

#### Connect panels

- [ ] solder male dupont (8 cm) and the - of both panels together

- [ ] solder male dupont (3 cm?) and the + of both panels together

### Supercap

- [ ] connect supercaps together in series. (GND on one to + on the other)

- [ ] solder the connected leads together

- [ ] wrap the 2 caps together using krylon tape

- [ ] connect Shottky diode to + on supercaps with ring towards the caps

- [ ] solder the diode and caps together

- [ ] solder female dupont (3 cm) to - on supercaps

- [ ] solder female dupont (3 cm) to diode input on supercaps

- [ ] solder female dupont (3 cm) to + on supercaps

- [ ] carefully wrap the individual connectors together using krylon tape
