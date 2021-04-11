package raildevices

// A turnout or railroad switch is a rail device used for changing the direction of a train to a diverging route.
// The difference to a normal switch is the output must not be permanent set to on, but only for a time period of 0.25-1s
// Second difference: A standard model turnout needs 2 physical outputs (left, right).

import (
	"github.com/gen2thomas/gobrail/internal/boardpin"
	"time"
)

const maxTime = time.Duration(time.Second)

// TurnoutDevice is describes a turnout
type TurnoutDevice struct {
	cmnOutDev    *CommonOutputDevice
	outputBranch *boardpin.Output
	outputMain   *boardpin.Output
}

// NewTurnout creates an instance of a turnout
func NewTurnout(co *CommonOutputDevice, outputBranch *boardpin.Output, outputMain *boardpin.Output) (s *TurnoutDevice) {
	s = &TurnoutDevice{
		cmnOutDev:    co,
		outputBranch: outputBranch,
		outputMain:   outputMain,
	}
	return
}

// SwitchOn will try to switch the turnout to diverging route
func (s *TurnoutDevice) SwitchOn() (err error) {
	if err = s.outputBranch.WriteValue(1); err != nil {
		return
	}
	s.cmnOutDev.TimingForStart()
	if err = s.outputBranch.WriteValue(0); err != nil {
		return
	}
	s.cmnOutDev.SetState(true)
	return
}

// SwitchOff will switch the turnout to main route
func (s *TurnoutDevice) SwitchOff() (err error) {
	if err = s.outputMain.WriteValue(1); err != nil {
		return
	}
	s.cmnOutDev.TimingForStop()
	if err = s.outputMain.WriteValue(0); err != nil {
		return
	}
	s.cmnOutDev.SetState(false)
	return
}

// StateChanged states true when turnout status was changed since last visit
func (s *TurnoutDevice) StateChanged(visitor string) (hasChanged bool, err error) {
	return s.cmnOutDev.StateChanged(visitor)
}

// IsOn means the track switch is switched to this direction, that the train will run
// to the inner circle or diverging route from the main route
//
// =CHOO-CHOO>====== --> IsOn = false (train runs the main route or outer circle)
//              \\
//               ||  --> IsOn = true  (train runs the inner circle)
//              //
func (s *TurnoutDevice) IsOn() bool {
	return s.cmnOutDev.IsOn()
}

// RailDeviceName gets the name of the turnout common output
func (s *TurnoutDevice) RailDeviceName() string {
	return s.cmnOutDev.RailDeviceName()
}

// Connect is connecting an input for use in Run()
func (s *TurnoutDevice) Connect(inputDevice Inputer) (err error) {
	return s.cmnOutDev.Connect(inputDevice)
}

// ConnectInverse is connecting an input for use in Run(), but with inversed action
func (s *TurnoutDevice) ConnectInverse(inputDevice Inputer) (err error) {
	return s.cmnOutDev.ConnectInverse(inputDevice)
}

// Run is called in a loop and will make action dependant on the input device
func (s *TurnoutDevice) Run() (err error) {
	return s.cmnOutDev.Run(s.SwitchOn, s.SwitchOff)
}

// ReleaseInput is used to unmap
func (s *TurnoutDevice) ReleaseInput() {
	s.cmnOutDev.ReleaseInput()
}
