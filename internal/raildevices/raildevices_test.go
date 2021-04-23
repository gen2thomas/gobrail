package raildevices

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

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

type timingTest struct {
	name string
	inT  Timing
	max  time.Duration
	expT Timing
}

func TestTimingLimit(t *testing.T) {
	// arrange
	assert := assert.New(t)
	var timingTests = []timingTest{
		{
			name: "unshrinked",
			inT:  Timing{Starting: time.Second, Stopping: time.Millisecond},
			max:  time.Duration(2 * time.Second),
			expT: Timing{Starting: time.Second, Stopping: time.Millisecond},
		},
		{
			name: "start limited",
			inT:  Timing{Starting: time.Second, Stopping: time.Millisecond},
			max:  time.Millisecond,
			expT: Timing{Starting: time.Millisecond, Stopping: time.Millisecond},
		},
		{
			name: "stop limited",
			inT:  Timing{Starting: time.Millisecond, Stopping: time.Second},
			max:  time.Duration(2 * time.Millisecond),
			expT: Timing{Starting: time.Millisecond, Stopping: time.Duration(2 * time.Millisecond)},
		},
	}
	for _, tt := range timingTests {
		t.Run(tt.name, func(t *testing.T) {
			// act
			tt.inT.Limit(tt.max)
			// assert
			assert.Equal(tt.expT, tt.inT)
		})
	}
}
