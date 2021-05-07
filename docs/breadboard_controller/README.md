## Using breadboards to get in touch

Playing around with breadboards is very nice. It reduces the need of soldering to a minimum or zero (when using pinout boards for terminals).
To startup quickly there are some breadboard schematics in the documentation folders. So buy some breadboards and startup prototyping your model railroad. Some examples can be found in folder "docs/breadboard_controller".

#### Notes for all example breadboards

The board named "Turnout-Simulation Board" is just to show the basic internal wiring of a magnetic turnout drive and is not needed in reality.

Address pins A0-A5 are "H" by internal pull-up resistors. Feel free to add wires from pins 1, 2, 3, 9, 11, 12 to GND for adjust the I2C address to your needs.

**Attention:** Protect your devices by removing the main power plug (18-20V=) during startup and before shutdown. Disregard may lead to destruction of your devices, especially your magnetic coils.

## Comparison I2C and proprietary hardware
#### I2C

Advantage
* old and well known standard
* with small improvements suitable for long range and high speed (fast mode ... ultra fast mode)
* large range of chips (IO's, PWM's, sensors for temperature, humidity or compass etc.)

Disadvantages
* software must be developed according to the used chip
* board must be developed or at least hardware must be soldered on a prototype board

#### Proprietary hardware

Advantage
* stable, well tested
* compatible with each other and following the common standards

Disadvantage
* closed source hardware and software
* custom components
* very expensive
* some principles seems to be very old (S88)

## What does that I2C board cost?
#### Basic "Type 2" board for drive 2 turnouts or other magnetic devices, using 4 hardware keys as inputs
* 8 channel (4I/4O) Type 2 board with PCA9501 (chip, 4 x IRLZ34, diodes, resistors, jumpers, prototype board) ca. 7 EUR

#### Operate 2 turnouts with active feedback to controller using the inputs
* Type 2 board, 4 opto-couplers, resistors ca. 8 EUR

#### Operate 4 LED's (or other load <= 20mA with active low) and 4 TTL inputs
* Type 2 board without MOSFET's ca. 4 EUR

#### Improve I2C signals for long range and/or pure wiring
* LTC43111 ca. 10 EUR

## Reducing costs and save environment by using old hardware
This is especially an ecological point.

* DIL opto-couplers can be found on a wide range of switching power supplies
* MOSFET transistors can be found on laptop main boards in large numbers, mostly SMD
* resistors (no SMD) are available on simple boards
* light barriers can be made from IR transmitters/receivers of remote controllers or using old fire detectors
* excellent 4 wire cables for I2C are USB cables (please use only shielded cables) or network CAT cables
* old USB or network connectors can be found on many old devices
