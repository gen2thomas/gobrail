# gobrail "A model railroad is just a gobot."
This project is based on the well structured library "gobot", see https://github.com/hybridgroup/gobot. 
Basically i2c bus devices are used here.

## Download/Install

The easiest way to install is to run `go get -u github.com/gen2thomas/gobrail`.

## Run without hardware & software modifications

The hardware of "Typ 2" is needed together with digispark, than simply `make run`.

## Run with another hardware

For initial tests the "Typ2" board can cleaned from MOSFETS's.
Using another adaptor is also possible. Please see the gobot documentation.

## State of development

This is prototype software. It is a long way to make it work in real model railroad environment.

## TODO's

* add next architecture levels (command interpreter)
* add sequences for output rail devices (signals, track switches)
* add sequences for input rail devices (toggle button)
* add first locomotive decoder (hopefully possible to gobot i2c devices)
