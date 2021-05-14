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

func TestNewBoardType2i(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	// act
	boardt2 := NewBoardType2i(new(adaptorMock), 0x05, "TestNewBoardType2i")
	// assert
	require.NotNil(boardt2)
	assert.Equal("TestNewBoardType2i", boardt2.name)
}

func TestNewBoardType2o(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	// act
	boardt2 := NewBoardType2o(new(adaptorMock), 0x05, "TestNewBoardType2o")
	// assert
	require.NotNil(boardt2)
	assert.Equal("TestNewBoardType2o", boardt2.name)
}

func TestNewBoardType2io(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	// act
	boardt2 := NewBoardType2io(new(adaptorMock), 0x05, "TestNewBoardType2io")
	// assert
	require.NotNil(boardt2)
	assert.Equal("TestNewBoardType2io", boardt2.name)
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
