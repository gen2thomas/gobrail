package raildevices

import (
	"github.com/gen2thomas/gobrail/internal/boardpin"
)

type ReadMock struct {
	callCounter int
	values      [5]uint8
	simError    error
}

type WriteMock struct {
	callCounter int
	values      [5]uint8
	simError    error
}

func NewInputMock(readMock *ReadMock) *boardpin.Input {
	return &boardpin.Input{
		ReadValue: func() (value uint8, err error) {
			return inputReadValueImpl(readMock)
		},
	}
}

func NewOutputMock(writeMock *WriteMock) *boardpin.Output {
	return &boardpin.Output{
		WriteValue: func(value uint8) (err error) {
			return ouputWriteValueImpl(writeMock, value)
		},
	}
}

func inputReadValueImpl(rm *ReadMock) (value uint8, err error) {
	rm.callCounter++
	return rm.values[rm.callCounter-1], rm.simError
}

func ouputWriteValueImpl(wm *WriteMock, value uint8) (err error) {
	wm.callCounter++
	wm.values[wm.callCounter-1] = value
	return wm.simError
}
