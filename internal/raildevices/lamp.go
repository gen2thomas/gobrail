package raildevices

// A lamp is a rail device used for
// simple lamps, neon-light simulation, blinking lamps

// LampDevice is describes a lamp
type LampDevice struct {
	cmnOutDev *CommonOutputDevice
}

// NewLamp creates an instance of a lamp
func NewLamp(boardsAPI BoardsAPIer, boardID string, boardPinNr uint8, railDeviceName string, timing Timing) (ld *LampDevice, err error) {
	var co *CommonOutputDevice
	if co, err = NewCommonOutput(boardsAPI, boardID, boardPinNr, railDeviceName, timing, "lamp"); err != nil {
		return
	}
	ld = &LampDevice{
		cmnOutDev: co,
	}
	return
}

// StateChanged states true when StdLamp status was changed since last visit
func (l *LampDevice) StateChanged(visitor string) (hasChanged bool, err error) {
	return l.cmnOutDev.StateChanged(visitor)
}

// IsOn states true when StdLamp is on
func (l *LampDevice) IsOn() bool {
	return l.cmnOutDev.IsOn()
}

// IsDefective states true when StdLamp is defective
func (l *LampDevice) IsDefective() bool {
	return l.cmnOutDev.IsDefective()
}

// SwitchOn will try to switch on the StdLamp
func (l *LampDevice) SwitchOn() (err error) {
	return l.cmnOutDev.SwitchOn()
}

// SwitchOff will switch off the StdLamp
func (l *LampDevice) SwitchOff() (err error) {
	return l.cmnOutDev.SwitchOff()
}

// MakeDefective causes the StdLamp in an simulated defective state
func (l *LampDevice) MakeDefective() (err error) {
	return l.cmnOutDev.MakeDefective()
}

// Repair will fix the simulated defective state
func (l *LampDevice) Repair() (err error) {
	return l.cmnOutDev.Repair()
}

// RailDeviceName gets the name of the lamp output
func (l *LampDevice) RailDeviceName() string {
	return l.cmnOutDev.RailDeviceName()
}

// Connect is connecting an input for use in Run()
func (l *LampDevice) Connect(inputDevice Inputer) (err error) {
	return l.cmnOutDev.Connect(inputDevice)
}

// ConnectInverse is connecting an input for use in Run(), but with inversed action
func (l *LampDevice) ConnectInverse(inputDevice Inputer) (err error) {
	return l.cmnOutDev.ConnectInverse(inputDevice)
}

// Run is called in a loop and will make action dependant on the input device
func (l *LampDevice) Run() (err error) {
	return l.cmnOutDev.Run(l.cmnOutDev.SwitchOn, l.cmnOutDev.SwitchOff)
}

// ReleaseInput is used to unmap
func (l *LampDevice) ReleaseInput() {
	l.cmnOutDev.ReleaseInput()
}
