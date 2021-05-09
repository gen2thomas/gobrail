package raildevicesapi

import (
	"fmt"
)

type runableDevice struct {
	Runner
	connectedInput Inputer
	inputInversion bool
	firstRun       bool
}

func newRunableDevice(outDev Runner) *runableDevice {
	return &runableDevice{
		Runner:   outDev,
		firstRun: true,
	}
}

// Connect is connecting an input for use in Run()
func (o *runableDevice) Connect(inputDevice Inputer, inputInversion bool) (err error) {
	if o.connectedInput != nil {
		return fmt.Errorf("The '%s' is already connected to an input '%s'", o.RailDeviceName(), o.connectedInput.RailDeviceName())
	}
	if o.RailDeviceName() == inputDevice.RailDeviceName() {
		return fmt.Errorf("Circular mapping blocked for '%s'", o.RailDeviceName())
	}
	o.connectedInput = inputDevice
	o.inputInversion = inputInversion
	return nil
}

// RunCommon is called in a loop and will make action, dependent on the input device
func (o *runableDevice) Run() (err error) {
	if o.connectedInput == nil {
		return fmt.Errorf("The '%s' can't run, please map to an input first", o.RailDeviceName())
	}
	var changed bool
	if changed, err = o.connectedInput.StateChanged(o.RailDeviceName()); err != nil {
		return err
	}
	fmt.Println("changed:", changed)
	if !(changed || o.firstRun) {
		return
	}
	o.firstRun = false
	if o.connectedInput.IsOn() != o.inputInversion {
		err = o.SwitchOn()
	} else {
		err = o.SwitchOff()
	}
	return
}

// ReleaseInput is used to unmap
func (o *runableDevice) ReleaseInput() {
	o.connectedInput = nil
}
