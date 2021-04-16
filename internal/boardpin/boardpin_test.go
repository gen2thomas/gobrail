package boardpin_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/gen2thomas/gobrail/internal/boardpin"
)

type pintypeTest struct {
	name          string
	pin           Pin
	typesList     []PinType
	mustContained bool
}

type pinNrTest struct {
	name     string
	numbers  PinNumbers
	expected []string
}

var allpinTypes = []PinType{
	Binary, BinaryR, BinaryW, NBinary, NBinaryR, NBinaryW,
	Analog, AnalogR, AnalogW,
	Memory, MemoryR, MemoryW,
}

func TestContainsPinType(t *testing.T) {
	// arrange
	assert := assert.New(t)
	var pintypeTests = []pintypeTest{
		{name: "Binary is in all pin types", pin: Pin{PinType: Binary}, typesList: allpinTypes, mustContained: true},
		{name: "BinaryR is in all pin types", pin: Pin{PinType: BinaryR}, typesList: allpinTypes, mustContained: true},
		{name: "BinaryW is in all pin types", pin: Pin{PinType: BinaryW}, typesList: allpinTypes, mustContained: true},
		{name: "NBinary is in all pin types", pin: Pin{PinType: NBinary}, typesList: allpinTypes, mustContained: true},
		{name: "NBinaryR is in all pin types", pin: Pin{PinType: NBinaryR}, typesList: allpinTypes, mustContained: true},
		{name: "NBinaryW is in all pin types", pin: Pin{PinType: NBinaryW}, typesList: allpinTypes, mustContained: true},
		{name: "Analog is in all pin types", pin: Pin{PinType: Analog}, typesList: allpinTypes, mustContained: true},
		{name: "AnalogR is in all pin types", pin: Pin{PinType: AnalogR}, typesList: allpinTypes, mustContained: true},
		{name: "AnalogW is in all pin types", pin: Pin{PinType: AnalogW}, typesList: allpinTypes, mustContained: true},
		{name: "Memory is in all pin types", pin: Pin{PinType: Memory}, typesList: allpinTypes, mustContained: true},
		{name: "MemoryR is in all pin types", pin: Pin{PinType: MemoryR}, typesList: allpinTypes, mustContained: true},
		{name: "MemoryW is in all pin types", pin: Pin{PinType: MemoryW}, typesList: allpinTypes, mustContained: true},
		{name: "Binary is not in Memory and BinaryR", pin: Pin{PinType: Binary}, typesList: []PinType{Memory, BinaryR}, mustContained: false},
		{name: "BinaryR is not in MemoryW", pin: Pin{PinType: BinaryR}, typesList: []PinType{MemoryW}, mustContained: false},
		{name: "BinaryR is not in empty list", pin: Pin{PinType: BinaryR}, typesList: []PinType{}, mustContained: false},
	}
	for _, pt := range pintypeTests {
		t.Run(pt.name, func(t *testing.T) {
			// act
			iscontained := pt.pin.PinTypeIsOneOf(pt.typesList)
			strMsg := pt.pin.PinType.String()
			// assert
			assert.Equal(pt.mustContained, iscontained)
			assert.NotContains(strMsg, "Unknown")
		})

	}
}

func TestUnknownPinTypeToString(t *testing.T) {
	// arrange
	assert := assert.New(t)
	// act
	strMsg := PinType(255).String()
	// assert
	assert.Contains(strMsg, "Unknown")
}

func TestPinNumbersToString(t *testing.T) {
	// arrange
	assert := assert.New(t)
	var pinNrTests = []pinNrTest{
		{
			name:     "One number",
			numbers:  PinNumbers{uint8(1): struct{}{}},
			expected: []string{"1"},
		},
		{
			name:     "Two numbers",
			numbers:  PinNumbers{uint8(1): struct{}{}, uint8(3): struct{}{}},
			expected: []string{"1", "3"},
		},
	}
	for _, pt := range pinNrTests {
		t.Run(pt.name, func(t *testing.T) {
			// act
			str := pt.numbers.String()
			// assert
			for _, estr := range pt.expected {
				assert.Contains(str, estr)
			}
		})
	}
}
