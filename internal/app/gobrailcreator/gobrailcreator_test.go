package gobrailcreator

type BoardsAPIMock struct {
	values              map[string]uint8
	callCounterBinMap   int
	callCounterAnaMap   int
	callCounterMemMap   int
	callCounterSetValue int
	callCounterGetValue int

	apiMapBinaryImpl func(boardID string, boardPinNr uint8, railDeviceName string) (err error)
	apiMapAnalogImpl func(boardID string, boardPinNr uint8, railDeviceName string) (err error)
	apiMapMemoryImpl func(boardID string, boardPinNrOrNegative int, railDeviceName string) (err error)
	apiSetValueImpl  func(railDeviceName string, value uint8) (err error)
	apiGetValueImpl  func(railDeviceName string) (value uint8, err error)
}

func NewBoardsAPIMock() *BoardsAPIMock {
	api := new(BoardsAPIMock)
	// values map simulates a board with mapped rail devices
	api.values = make(map[string]uint8)
	api.apiMapBinaryImpl = func(boardID string, boardPinNr uint8, railDeviceName string) (err error) {
		api.values[railDeviceName] = 0
		return
	}
	api.apiMapAnalogImpl = func(boardID string, boardPinNr uint8, railDeviceName string) (err error) {
		api.values[railDeviceName] = 0
		return
	}
	api.apiMapMemoryImpl = func(boardID string, boardPinNrOrNegative int, railDeviceName string) (err error) {
		api.values[railDeviceName] = 0
		return
	}
	api.apiSetValueImpl = func(railDeviceName string, value uint8) (err error) {
		api.values[railDeviceName] = value
		return
	}
	api.apiGetValueImpl = func(railDeviceName string) (value uint8, err error) {
		value = api.values[railDeviceName]
		return
	}

	return api
}

func (ba *BoardsAPIMock) MapBinaryPin(boardID string, boardPinNr uint8, railDeviceName string) (err error) {
	ba.callCounterBinMap++
	return ba.apiMapBinaryImpl(boardID, boardPinNr, railDeviceName)
}
func (ba *BoardsAPIMock) MapAnalogPin(boardID string, boardPinNr uint8, railDeviceName string) (err error) {
	ba.callCounterAnaMap++
	return ba.apiMapAnalogImpl(boardID, boardPinNr, railDeviceName)
}
func (ba *BoardsAPIMock) MapMemoryPin(boardID string, boardPinNrOrNegative int, railDeviceName string) (err error) {
	ba.callCounterMemMap++
	return ba.apiMapMemoryImpl(boardID, boardPinNrOrNegative, railDeviceName)
}
func (ba *BoardsAPIMock) GetValue(railDeviceName string) (value uint8, err error) {
	ba.callCounterGetValue++
	return ba.apiGetValueImpl(railDeviceName)
}
func (ba *BoardsAPIMock) SetValue(railDeviceName string, value uint8) (err error) {
	ba.callCounterSetValue++
	return ba.apiSetValueImpl(railDeviceName, value)
}
