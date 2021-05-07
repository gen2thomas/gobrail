# gobrail "A model railroad is just a gobot."
This project is based on the well structured library "gobot", see https://github.com/hybridgroup/gobot. 
Basically i2c bus devices are used here.

## Download/Install

The easiest way to install is to run `go get -u github.com/gen2thomas/gobrail`.

## Run without hardware & software modifications
#### Configuration

The hardware structure will be configured with a plan file in json format. See folder docs/breadboard_controller/ for examples.

#### amd64 host system

The hardware of "Type 2" is needed together with digispark, than simply `make run`.
Please have a look at folder docs/breadboard_controller/ for using with a breadboard.

#### arm/amd64 targets

First run `make` to create all binaries for target systems. Choose the binary for your target from output folder and copy to your target device.
* gobrail => amd64
* gobrail_raspi => Raspberry Pi

Make the file executable at your target device and start it by e.g. `./gobrail_raspi -help` to list command line parameters and its defaults.

## Hardware modifications

For initial tests the "Type 2" board can cleaned from MOSFETS's.
Using another adaptor is also possible. Please see the gobot documentation for this topic.

## State of development

This is prototype software. It is a long way to make it work in real model railroad environment.

#### Supported output rail devices

* lamp - single output on/off
* two light signal - two outputs green on -> red off and vice versa
* turnout - two outputs switched on, configurable between 0-1 second, to switch between main and branch

#### Supported input rail devices

* button - read one input
* toggle button - read one input, use rising edge to change the state

## TODO's

* improve timing by using events and/or concurrency
* add configuration interface
* add first locomotive decoder (hopefully possible to gobot devices)
* virtual boards (a button or lamp can be mapped to an virtual IO, which provides an "external service")
