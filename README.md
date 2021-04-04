# gobrail "A model railroad is just a gobot."
This project is based on the well structured library "gobot", see https://github.com/hybridgroup/gobot. 
Basically i2c bus devices are used here.

## Download/Install

The easiest way to install is to run `go get -u github.com/gen2thomas/gobrail`.

## Run without hardwrae & software modifications

The hardware of "Typ 2" is needed together with digispark, than simply `make run`.

## Run with another hardware

For initial tests the "Typ2" board can cleaned from MOSFETS's.
Using another adaptor is also possible. Please see the gobot documentation.

## State of development

This is prototype software. It is a long way to make it work in real model railroad environment.

## TODO's

* documentation for "board" and planned software architecture
* tests for "board"
* add negotiation possibility
* add next architecture levels (boardsapi, rail)
* add sequences for rail devices (lamps, signals, track switches)
* add first locomotive decoder (hopefully possible to gobot i2c devices)
