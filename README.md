# gobrail "A model railroad is just a gobot."
This project is based on the well structured library "gobot", see https://github.com/hybridgroup/gobot. 
Basically i2c bus devices are used here.

## Download/Install

The easiest way to install is to run `go get -u github.com/gen2thomas/gobrail`.

## Run without hardware & software modifications

The hardware of "Typ 2" is needed together with digispark, than simply `make run`.
Please have a look at docs/images/PCA9501_Lamps_Buttons.png for using with a breadboard.

## Run with another hardware

For initial tests the "Typ2" board can cleaned from MOSFETS's.
Using another adaptor is also possible. Please see the gobot documentation.

## State of development

This is prototype software. It is a long way to make it work in real model railroad environment.

### Supported output rail devices
* lamp - single output on/off
* two light signal - two outputs green on -> red off and vice versa
* turnout - two outputs on for 0-1sec to switch between main and branch

### Supported input rail devices
* button - read one input
* toggle button - read one input, use rising edge to change the state

### Instantiating and connecting devices
* possible with simple go programming knowledge before main loop, see examples

### Running
* possible with simple go programming knowledge in main loop, see examples

## TODO's

### General
* improve timing by using events and/or parallelism

### Hardware UI interaction
* add rail runner
* add configuration interface
* add first locomotive decoder (hopefully possible to gobot i2c devices)

### Software UI interaction
* virtual boards (a button or lamp can be mapped to an virtual IO, which provides an "external service")
