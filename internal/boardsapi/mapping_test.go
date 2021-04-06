package boardsapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gen2thomas/gobrail/internal/board"
)

func TestFindRailDevice(t *testing.T) {
	// arrange
	assert := assert.New(t)
	// boards
	b1 := board.NewBoardForTestWithoutChips("TestBoard1", 1, 1, 1)
	b2 := board.NewBoardForTestWithoutChips("TestBoard2", 1, 0, 0)
	bm := make(BoardsMap)
	bm["TestBoard1"] = b1
	bm["TestBoard2"] = b2
	// mapping
	mps := make(APIPinsMap)
	mps["k1"] = &apiPin{boardID: "TestBoard1", boardPinNr: 0, railDeviceName: "bin device 1"}
	mps["k2"] = &apiPin{boardID: "TestBoard1", boardPinNr: 1, railDeviceName: "ana device 1"}
	mps["k3"] = &apiPin{boardID: "TestBoard1", boardPinNr: 2, railDeviceName: "mem device 1"}
	mps["k4"] = &apiPin{boardID: "TestBoard2", boardPinNr: 0, railDeviceName: "bin device 2"}
	api := &BoardsAPI{
		mappedPins: mps,
		boards:     bm,
	}
	// act
	r1bin := api.FindRailDevice("TestBoard1", 0)
	r1ana := api.FindRailDevice("TestBoard1", 1)
	r1mem := api.FindRailDevice("TestBoard1", 2)
	r1no := api.FindRailDevice("TestBoard1", 3)
	r2bin := api.FindRailDevice("TestBoard2", 0)
	r3no := api.FindRailDevice("TestBoard3", 0)
	// assert
	assert.Equal("k1", r1bin)
	assert.Equal("k2", r1ana)
	assert.Equal("k3", r1mem)
	assert.Equal("", r1no)
	assert.Equal("k4", r2bin)
	assert.Equal("", r3no)
}

func TestFindRailDeviceWithoutMappedPinsGetsEmptyString(t *testing.T) {
	// arrange
	assert := assert.New(t)
	// boards
	b1 := board.NewBoardForTestWithoutChips("TestBoard1", 1, 1, 1)
	bm := make(BoardsMap)
	bm["TestBoard1"] = b1
	// mapping
	mps := make(APIPinsMap)
	api := &BoardsAPI{
		mappedPins: mps,
		boards:     bm,
	}
	// act
	r1no := api.FindRailDevice("TestBoard1", 0)
	// assert
	assert.Equal("", r1no)
}

func TestGetFreeAPIPinsWithoutBoardGetsEmptyList(t *testing.T) {
	// arrange
	assert := assert.New(t)
	api := &BoardsAPI{
		mappedPins: make(APIPinsMap),
		boards:     make(BoardsMap),
	}
	// act
	fp := api.GetFreeAPIPins("NoExistend", board.Binary)
	// assert
	assert.Equal(0, len(fp))
}

func TestGetFreeAPIPins(t *testing.T) {
	// arrange
	assert := assert.New(t)
	b1 := board.NewBoardForTestWithoutChips("TestBoard1", 2, 5, 1)
	bm := make(BoardsMap)
	bm["TestBoard1"] = b1
	api := &BoardsAPI{
		mappedPins: make(APIPinsMap),
		boards:     bm,
	}
	// act
	fp1bin := api.GetFreeAPIPins("TestBoard1", board.Binary)
	fp1ana := api.GetFreeAPIPins("TestBoard1", board.Analog)
	fp1mem := api.GetFreeAPIPins("TestBoard1", board.Memory)
	// assert
	assert.Equal(2, len(fp1bin))
	assert.Equal(5, len(fp1ana))
	assert.Equal(1, len(fp1mem))
}

func TestGetMappedAPIPinsWithoutBoardGetsEmptyList(t *testing.T) {
	// arrange
	assert := assert.New(t)
	api := &BoardsAPI{
		mappedPins: make(APIPinsMap),
		boards:     make(BoardsMap),
	}
	// act
	mp := api.GetMappedAPIPins("NoExistend", board.Binary)
	// assert
	assert.Equal(0, len(mp))
}

func TestGetMappedAPIPins(t *testing.T) {
	// arrange
	assert := assert.New(t)
	b1 := board.NewBoardForTestWithoutChips("TestBoard1", 3, 3, 3)
	bm := make(BoardsMap)
	bm["TestBoard1"] = b1
	mps := make(APIPinsMap)
	mps["k1"] = &apiPin{boardID: "TestBoard1", boardPinNr: 4, railDeviceName: "ana device 1"}
	mps["k2"] = &apiPin{boardID: "TestBoard1", boardPinNr: 1, railDeviceName: "bin device 1"}
	mps["k3"] = &apiPin{boardID: "TestBoard1", boardPinNr: 5, railDeviceName: "ana device 2"}
	mps["k4"] = &apiPin{boardID: "TestBoard1", boardPinNr: 8, railDeviceName: "mem device 1"}
	mps["k5"] = &apiPin{boardID: "TestBoard1", boardPinNr: 0, railDeviceName: "bin device 2"}
	mps["k6"] = &apiPin{boardID: "TestBoard1", boardPinNr: 2, railDeviceName: "bin device 3"}
	api := &BoardsAPI{
		mappedPins: mps,
		boards:     bm,
	}
	// act
	mp1bin := api.GetMappedAPIPins("TestBoard1", board.Binary)
	mp1ana := api.GetMappedAPIPins("TestBoard1", board.Analog)
	mp1mem := api.GetMappedAPIPins("TestBoard1", board.Memory)
	// assert
	assert.Equal(3, len(mp1bin))
	assert.Equal(2, len(mp1ana))
	assert.Equal(1, len(mp1mem))
}

func TestMapPin(t *testing.T) {
	// arrange
	assert := assert.New(t)
	b1 := board.NewBoardForTestWithoutChips("TestBoard1", 1, 1, 1)
	bm := make(BoardsMap)
	bm["TestBoard1"] = b1
	api := &BoardsAPI{
		mappedPins: make(APIPinsMap),
		boards:     bm,
	}
	// act
	err1 := api.MapPin("TestBoard1", 1, "ana device 1")
	err2 := api.MapPin("TestBoard1", 0, "bin device 1")
	err3 := api.MapPin("TestBoard1", 2, "mem device 1")
	// no access to mapped pins for testing purposes other than calling a function
	// GetMappedAPIPins() would be also possible
	k1mem := api.FindRailDevice("TestBoard1", 2)
	k1ana := api.FindRailDevice("TestBoard1", 1)
	k1bin := api.FindRailDevice("TestBoard1", 0)
	// assert
	assert.Nil(err1)
	assert.Nil(err2)
	assert.Nil(err3)
	assert.Equal("bin_device_1", k1bin)
	assert.Equal("ana_device_1", k1ana)
	assert.Equal("mem_device_1", k1mem)
}

func TestMapPinWithAlreadyMappedRailDeviceGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	b1 := board.NewBoardForTestWithoutChips("TestBoard1", 2, 0, 0)
	bm := make(BoardsMap)
	bm["TestBoard1"] = b1
	// mapped pins
	mps := make(APIPinsMap)
	mps["bin_device_1"] = &apiPin{boardID: "TestBoard1", boardPinNr: 0, railDeviceName: "bin device 1"}
	api := &BoardsAPI{
		mappedPins: mps,
		boards:     bm,
	}
	// act
	err := api.MapPin("TestBoard1", 1, "bin device 1")
	// assert
	require.NotNil(err)
	assert.Contains(err.Error(), "Rail device")
}

func TestMapPinWithAlreadyMappedBoardPinGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	b1 := board.NewBoardForTestWithoutChips("TestBoard1", 2, 0, 0)
	bm := make(BoardsMap)
	bm["TestBoard1"] = b1
	// mapped pins
	mps := make(APIPinsMap)
	mps["k1"] = &apiPin{boardID: "TestBoard1", boardPinNr: 0, railDeviceName: "bin device 1"}
	api := &BoardsAPI{
		mappedPins: mps,
		boards:     bm,
	}
	// act
	err := api.MapPin("TestBoard1", 0, "bin device 2")
	// assert
	require.NotNil(err)
	assert.Contains(err.Error(), "Pin already")
}

func TestMapPinNextFree(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	b1 := board.NewBoardForTestWithoutChips("TestBoard1", 2, 2, 2)
	bm := make(BoardsMap)
	bm["TestBoard1"] = b1
	// mapped pins
	mps := make(APIPinsMap)
	mps["k1"] = &apiPin{boardID: "TestBoard1", boardPinNr: 0, railDeviceName: "bin device 1"}
	mps["k2"] = &apiPin{boardID: "TestBoard1", boardPinNr: 2, railDeviceName: "ana device 1"}
	mps["k3"] = &apiPin{boardID: "TestBoard1", boardPinNr: 4, railDeviceName: "mem device 1"}
	api := &BoardsAPI{
		mappedPins: mps,
		boards:     bm,
	}
	// act
	err1 := api.MapPinNextFree("TestBoard1", board.Binary, "bin device 2")
	err2 := api.MapPinNextFree("TestBoard1", board.Analog, "ana device 2")
	err3 := api.MapPinNextFree("TestBoard1", board.Memory, "mem device 2")
	// no access to mapped pins for testing purposes other than calling a function
	// FindRailDevice() would be also possible
	mp1bin := api.GetMappedAPIPins("TestBoard1", board.Binary)
	mp1ana := api.GetMappedAPIPins("TestBoard1", board.Analog)
	mp1mem := api.GetMappedAPIPins("TestBoard1", board.Memory)
	// assert
	require.Nil(err1)
	require.Nil(err2)
	require.Nil(err3)
	assert.Equal(2, len(mp1bin))
	assert.Equal(2, len(mp1ana))
	assert.Equal(2, len(mp1mem))
	assert.Contains(mp1bin.String(), "bin device 2")
	assert.Contains(mp1ana.String(), "ana device 2")
	assert.Contains(mp1mem.String(), "mem device 2")
}

func TestMapPinNextFreeWithoutFreePinsGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	b1 := board.NewBoardForTestWithoutChips("TestBoard1", 1, 1, 1)
	bm := make(BoardsMap)
	bm["TestBoard1"] = b1
	// mapped pins
	mps := make(APIPinsMap)
	mps["k1"] = &apiPin{boardID: "TestBoard1", boardPinNr: 0, railDeviceName: "bin device 1"}
	mps["k2"] = &apiPin{boardID: "TestBoard1", boardPinNr: 1, railDeviceName: "ana device 1"}
	mps["k3"] = &apiPin{boardID: "TestBoard1", boardPinNr: 2, railDeviceName: "mem device 1"}
	api := &BoardsAPI{
		mappedPins: mps,
		boards:     bm,
	}
	// act
	err1bin := api.MapPinNextFree("TestBoard1", board.Binary, "bin device 2")
	err1ana := api.MapPinNextFree("TestBoard1", board.Analog, "ana device 2")
	err1mem := api.MapPinNextFree("TestBoard1", board.Memory, "mem device 2")
	// assert
	require.NotNil(err1bin)
	require.NotNil(err1ana)
	require.NotNil(err1mem)
	assert.Contains(err1bin.Error(), "No free pin at")
	assert.Contains(err1ana.Error(), "No free pin at")
	assert.Contains(err1mem.Error(), "No free pin at")
	assert.Contains(err1bin.Error(), "bin device 2")
	assert.Contains(err1ana.Error(), "ana device 2")
	assert.Contains(err1mem.Error(), "mem device 2")
}

func TestReleasePin(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	b1 := board.NewBoardForTestWithoutChips("TestBoard1", 1, 0, 0)
	bm := make(BoardsMap)
	bm["TestBoard1"] = b1
	// mapped pins
	mps := make(APIPinsMap)
	mps["bin_device_1"] = &apiPin{boardID: "TestBoard1", boardPinNr: 0, railDeviceName: "bin device 1"}
	api := &BoardsAPI{
		mappedPins: mps,
		boards:     bm,
	}
	// act
	err := api.ReleasePin("bin device 1")
	// no access to mapped pins for testing purposes other than calling a function
	mp1bin := api.GetMappedAPIPins("TestBoard1", board.Binary)
	// assert
	require.Nil(err)
	assert.Equal(0, len(mp1bin))
}

func TestReleasePinOfNotMappedRailDeviceGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	b1 := board.NewBoardForTestWithoutChips("TestBoard1", 1, 0, 0)
	bm := make(BoardsMap)
	bm["TestBoard1"] = b1
	// mapped pins
	mps := make(APIPinsMap)
	mps["bin_device_1"] = &apiPin{boardID: "TestBoard1", boardPinNr: 0, railDeviceName: "bin device 1"}
	api := &BoardsAPI{
		mappedPins: mps,
		boards:     bm,
	}
	// act
	err := api.ReleasePin("bin device 2")
	// assert
	require.NotNil(err)
	assert.Contains(err.Error(), "not mapped")
	assert.Contains(err.Error(), "bin device 2")
}

func TestReleasePinWithoutMappedPinsGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	b1 := board.NewBoardForTestWithoutChips("TestBoard1", 1, 0, 0)
	bm := make(BoardsMap)
	bm["TestBoard1"] = b1
	api := &BoardsAPI{
		mappedPins: make(APIPinsMap),
		boards:     bm,
	}
	// act
	err := api.ReleasePin("a device")
	// assert
	require.NotNil(err)
	assert.Contains(err.Error(), "not mapped")
	assert.Contains(err.Error(), "a device")
}
