package raildevices

type BoardsAPIMock struct {
	apiMapBinaryImpl func(boardID string, boardPinNr uint8, railDeviceName string) (err error)
	apiMapAnalogImpl func(boardID string, boardPinNr uint8, railDeviceName string) (err error)
	apiMapMemoryImpl func(boardID string, boardPinNrOrNegative int, railDeviceName string) (err error)
	apiSetValueImpl  func(railDeviceName string, value uint8) (err error)
	apiGetValueImpl  func(railDeviceName string) (value uint8, err error)
}

func (ba *BoardsAPIMock) MapBinaryPin(boardID string, boardPinNr uint8, railDeviceName string) (err error) {
	return ba.apiMapBinaryImpl(boardID, boardPinNr, railDeviceName)
}
func (ba *BoardsAPIMock) MapAnalogPin(boardID string, boardPinNr uint8, railDeviceName string) (err error) {
	return ba.apiMapAnalogImpl(boardID, boardPinNr, railDeviceName)
}
func (ba *BoardsAPIMock) MapMemoryPin(boardID string, boardPinNrOrNegative int, railDeviceName string) (err error) {
	return ba.apiMapMemoryImpl(boardID, boardPinNrOrNegative, railDeviceName)
}
func (ba *BoardsAPIMock) GetValue(railDeviceName string) (value uint8, err error) {
	return ba.apiGetValueImpl(railDeviceName)
}
func (ba *BoardsAPIMock) SetValue(railDeviceName string, value uint8) (err error) {
	return ba.apiSetValueImpl(railDeviceName, value)
}
