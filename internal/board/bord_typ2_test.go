package board

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gobot.io/x/gobot/drivers/i2c"

	"github.com/gen2thomas/gobrail/internal/boardpin"
)

type adaptorMock struct {
	name string
}

func TestNewBoardTyp2(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	// act
	boardt2 := NewBoardTyp2(new(adaptorMock), 0x05, "TestNewBoardTyp2")
	// assert
	require.NotNil(boardt2)
	assert.Equal("TestNewBoardTyp2", boardt2.name)
}

func TestWriteGPIOWithoutDriverFails(t *testing.T) {
	// arrange
	assert := assert.New(t)
	boardt2 := &Board{}
	// act
	err := boardt2.writeGPIO(&boardpin.Pin{}, 2)
	// assert
	assert.NotNil(err)
}

func TestReadGPIOWithoutDriverFails(t *testing.T) {
	// arrange
	assert := assert.New(t)
	boardt2 := &Board{}
	// act
	_, err := boardt2.readGPIO(&boardpin.Pin{})
	// assert
	assert.NotNil(err)
}

func TestWriteEEPROMWithoutDriverFails(t *testing.T) {
	// arrange
	assert := assert.New(t)
	boardt2 := &Board{}
	// act
	err := boardt2.writeEEPROM(&boardpin.Pin{}, 1)
	// assert
	assert.NotNil(err)
}

func TestReadEEPROMWithoutDriverFails(t *testing.T) {
	// arrange
	assert := assert.New(t)
	boardt2 := &Board{}
	// act
	_, err := boardt2.readEEPROM(&boardpin.Pin{})
	// assert
	assert.NotNil(err)
}

func (a *adaptorMock) GetConnection(address int, bus int) (device i2c.Connection, err error) { return }
func (a *adaptorMock) GetDefaultBus() int                                                    { return 0 }
