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
	*CommonOutputDevice
	outputBranch *boardpin.Output
	outputMain   *boardpin.Output
}

// NewTurnout creates an instance of a turnout
func NewTurnout(co *CommonOutputDevice, outputBranch *boardpin.Output, outputMain *boardpin.Output) (s *TurnoutDevice) {
	s = &TurnoutDevice{
		CommonOutputDevice: co,
		outputBranch:       outputBranch,
		outputMain:         outputMain,
	}
	return
}

// SwitchOn will try to switch the turnout to diverging route
// from the main route or the inner circle
//
// =CHOO-CHOO>====== --> IsOn = false (train runs the main route or outer circle)
//              \\
//               ||  --> IsOn = true  (train runs the inner circle)
//              //
func (s *TurnoutDevice) SwitchOn() (err error) {
	if err = s.outputBranch.WriteValue(1); err != nil {
		return
	}
	s.TimingForStart()
	if err = s.outputBranch.WriteValue(0); err != nil {
		return
	}
	s.SetState(true)
	return
}

// SwitchOff will switch the turnout to main route
func (s *TurnoutDevice) SwitchOff() (err error) {
	if err = s.outputMain.WriteValue(1); err != nil {
		return
	}
	s.TimingForStop()
	if err = s.outputMain.WriteValue(0); err != nil {
		return
	}
	s.SetState(false)
	return
}

// Run is called in a loop and will make action dependant on the input device
func (s *TurnoutDevice) Run() (err error) {
	return s.RunCommon(s.SwitchOn, s.SwitchOff)
}
