package board

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPinsOfType(t *testing.T) {
	// arrange
	assert := assert.New(t)
	boardPins := boardPinsMap{
		0: {pinType: Binary},
		1: {pinType: Memory},
		2: {pinType: Analog},
		3: {pinType: Analog},
		4: {pinType: Binary},
		5: {pinType: Binary},
	}

	testBoard := &Board{
		name: "TestPinsOfType",
		pins: boardPins,
	}

	// act
	pinsBin := testBoard.PinsOfType(Binary)
	pinsAna := testBoard.PinsOfType(Analog)
	pinsMem := testBoard.PinsOfType(Memory)

	// assert
	assert.Equal(len(pinsBin), 3)
	assert.Equal(len(pinsAna), 2)
	assert.Equal(len(pinsMem), 1)
}
