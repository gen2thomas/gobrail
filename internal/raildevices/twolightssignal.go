package raildevices

// A two light signal is a rail device used for sign "pass" or "stop" with lamps (commonly in colours "green"/"red")
// Both lights can't be set at the same time.
// The output is static set in difference to a semaphore signal, which is more like a "turnout" to contros.

// TwoLightSignalDevice is describes a signal with two lights
type TwoLightSignalDevice struct {
	railDeviceNameStop string
	cmnOutDev          *CommonOutputDevice
}

// NewTwoLightSignal creates an instance of a light signal with two lights
func NewTwoLightSignal(boardsAPI BoardsAPIer, boardID string, boardPinNr uint8, railDeviceName string, boardPinNrStop uint8, timing Timing) (s *TwoLightSignalDevice, err error) {
	var co *CommonOutputDevice
	if co, err = NewCommonOutput(boardsAPI, boardID, boardPinNr, railDeviceName, timing, "two light signal"); err != nil {
		return
	}
	railDevicerailDeviceNameStop := railDeviceName + " stop"
	if err = boardsAPI.MapBinaryPin(boardID, boardPinNrStop, railDevicerailDeviceNameStop); err != nil {
		return
	}
	s = &TwoLightSignalDevice{
		railDeviceNameStop: railDevicerailDeviceNameStop,
		cmnOutDev:          co,
	}
	return
}

// SwitchOn will try to switch off the "stop" light (e.g. red colour)
// and immediatally switch on the "can pass" light (e.g. green colour)
func (s *TwoLightSignalDevice) SwitchOn() (err error) {
	s.cmnOutDev.TimingForStart()
	if err = s.cmnOutDev.BoardsAPI.SetValue(s.railDeviceNameStop, 0); err != nil {
		return
	}
	if err = s.cmnOutDev.BoardsAPI.SetValue(s.RailDeviceName(), 1); err != nil {
		return
	}
	s.cmnOutDev.SetState(true)
	return
}

// SwitchOff will try to switch off the "can pass" light (e.g. green colour)
// and immediatally switch on the "stop" light (e.g. red colour)
func (s *TwoLightSignalDevice) SwitchOff() (err error) {
	s.cmnOutDev.TimingForStop()
	if err = s.cmnOutDev.BoardsAPI.SetValue(s.RailDeviceName(), 0); err != nil {
		return
	}
	if err = s.cmnOutDev.BoardsAPI.SetValue(s.railDeviceNameStop, 1); err != nil {
		return
	}
	s.cmnOutDev.SetState(false)
	return
}

// StateChanged states true when light signal status was changed since last visit
func (s *TwoLightSignalDevice) StateChanged(visitor string) (hasChanged bool, err error) {
	return s.cmnOutDev.StateChanged(visitor)
}

// IsOn means the "can pass" position (e.g. green colour)
func (s *TwoLightSignalDevice) IsOn() bool {
	return s.cmnOutDev.IsOn()
}

// RailDeviceName gets the name of the light signal common output
func (s *TwoLightSignalDevice) RailDeviceName() string {
	return s.cmnOutDev.RailDeviceName()
}

// Connect is connecting an input for use in Run()
func (s *TwoLightSignalDevice) Connect(inputDevice Inputer) (err error) {
	return s.cmnOutDev.Connect(inputDevice)
}

// ConnectInverse is connecting an input for use in Run(), but with inversed action
func (s *TwoLightSignalDevice) ConnectInverse(inputDevice Inputer) (err error) {
	return s.cmnOutDev.ConnectInverse(inputDevice)
}

// Run is called in a loop and will make action dependant on the input device
func (s *TwoLightSignalDevice) Run() (err error) {
	return s.cmnOutDev.Run(s.SwitchOn, s.SwitchOff)
}

// ReleaseInput is used to unmap
func (s *TwoLightSignalDevice) ReleaseInput() {
	s.cmnOutDev.ReleaseInput()
}
