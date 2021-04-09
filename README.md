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

## TODO's

### Hardware UI interaction
* add sequences for output rail devices (signal, railroad switch)
* interaction between turnout and signal ("if left than green, else red" or vice versa)
* add next architecture levels (command interpreter, rail runner)
* add first locomotive decoder (hopefully possible to gobot i2c devices)

### Software UI interaction
* virtual boards (a button or lamp can be mapped to an virtual IO, which provides an "external service")

### General
* improve timing by using events
