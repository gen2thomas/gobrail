package board

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gobot.io/x/gobot/drivers/i2c"
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
	assert.Equal(boardt2.name, "TestNewBoardTyp2")
}

func (a *adaptorMock) GetConnection(address int, bus int) (device i2c.Connection, err error) { return }
func (a *adaptorMock) GetDefaultBus() int                                                    { return 0 }
