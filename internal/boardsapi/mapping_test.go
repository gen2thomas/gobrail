package boardsapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot"
)

type boards1Mock struct {
	name    string
	binPins uint8
	anaPins uint8
	memPins uint8
}

func TestFindRailDevice(t *testing.T) {
	// arrange
	assert := assert.New(t)
	// boards
	bm := make(BoardsMap)
	bm["TestBoard1"] = &boards1Mock{name: "TestBoard1"}
	bm["TestBoard2"] = &boards1Mock{name: "TestBoard2"}
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
	bm := make(BoardsMap)
	bm["TestBoard1"] = &boards1Mock{name: "TestBoard1"}
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
	fpb := api.GetFreeAPIBinaryPins("NoExistend")
	fpa := api.GetFreeAPIAnalogPins("NoExistend")
	fpm := api.GetFreeAPIMemoryPins("NoExistend")
	// assert
	assert.Equal(0, len(fpb))
	assert.Equal(0, len(fpa))
	assert.Equal(0, len(fpm))
}

func TestGetFreeAPIPins(t *testing.T) {
	// arrange
	assert := assert.New(t)
	bm := make(BoardsMap)
	bm["TestBoard1"] = &boards1Mock{name: "TestBoard1", binPins: 2, anaPins: 5, memPins: 1}
	api := &BoardsAPI{
		mappedPins: make(APIPinsMap),
		boards:     bm,
	}
	// act
	fp1bin := api.GetFreeAPIBinaryPins("TestBoard1")
	fp1ana := api.GetFreeAPIAnalogPins("TestBoard1")
	fp1mem := api.GetFreeAPIMemoryPins("TestBoard1")
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
	mpb := api.GetMappedAPIBinaryPins("NoExistend")
	mpa := api.GetMappedAPIAnalogPins("NoExistend")
	mpm := api.GetMappedAPIMemoryPins("NoExistend")
	// assert
	assert.Equal(0, len(mpb))
	assert.Equal(0, len(mpa))
	assert.Equal(0, len(mpm))
}

func TestGetMappedAPIPins(t *testing.T) {
	// arrange
	assert := assert.New(t)
	bm := make(BoardsMap)
	bm["TestBoard1"] = &boards1Mock{name: "TestBoard1", binPins: 3, anaPins: 3, memPins: 3}
	mps := make(APIPinsMap)
	mps["k1"] = &apiPin{boardID: "TestBoard1", boardPinNr: 4, railDeviceName: "device pin4"}
	mps["k2"] = &apiPin{boardID: "TestBoard1", boardPinNr: 1, railDeviceName: "device pin1"}
	mps["k3"] = &apiPin{boardID: "TestBoard1", boardPinNr: 5, railDeviceName: "device pin5"}
	mps["k4"] = &apiPin{boardID: "TestBoard1", boardPinNr: 8, railDeviceName: "device pin8"}
	mps["k5"] = &apiPin{boardID: "TestBoard1", boardPinNr: 0, railDeviceName: "device pin0"}
	mps["k6"] = &apiPin{boardID: "TestBoard1", boardPinNr: 2, railDeviceName: "device pin2"}
	api := &BoardsAPI{
		mappedPins: mps,
		boards:     bm,
	}
	// act
	mp := api.getMappedAPIPins("TestBoard1", func() map[uint8]struct{} { return createBoardPinNumbersMap(1, 4) })
	// assert
	assert.Equal(3, len(mp))
	assert.Contains(mp.String(), "device pin1")
	assert.Contains(mp.String(), "device pin2")
	assert.Contains(mp.String(), "device pin4")
}

func TestMapPin(t *testing.T) {
	// arrange
	assert := assert.New(t)
	bm := make(BoardsMap)
	bm["TestBoard1"] = &boards1Mock{name: "TestBoard1", binPins: 1, anaPins: 1, memPins: 1}
	api := &BoardsAPI{
		mappedPins: make(APIPinsMap),
		boards:     bm,
	}
	// act
	err := api.mapPin("TestBoard1", 1, "device 1",
		func(string) (freePins APIPinsMap) {
			return APIPinsMap{"freekey": &apiPin{boardID: "TestBoard1", boardPinNr: 1}}
		})
	// assert
	assert.Nil(err)
	assert.Equal(1, len(api.mappedPins))
	assert.Equal("device 1", api.mappedPins["device_1"].railDeviceName)
}

func TestMapPinWhenNegativeBoardPinThanMapNextFree(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	bm := make(BoardsMap)
	bm["TestBoard1"] = &boards1Mock{name: "TestBoard1"}
	api := &BoardsAPI{
		mappedPins: make(APIPinsMap),
		boards:     bm,
	}
	// act
	err := api.mapPin("TestBoard1", -1, "device 1",
		func(string) (freePins APIPinsMap) {
			return APIPinsMap{"freekey": &apiPin{boardID: "TestBoard1", boardPinNr: 1}}
		})
	// assert
	require.Nil(err)
	assert.Equal(1, len(api.mappedPins))
	assert.Equal("device 1", api.mappedPins["device_1"].railDeviceName)
}

func TestMapPinWithAlreadyMappedRailDeviceGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	bm := make(BoardsMap)
	bm["TestBoard1"] = &boards1Mock{name: "TestBoard1"}
	// mapped pins
	mps := make(APIPinsMap)
	mps["device_1"] = &apiPin{boardID: "TestBoard1", boardPinNr: 0, railDeviceName: "device 1"}
	api := &BoardsAPI{
		mappedPins: mps,
		boards:     bm,
	}
	// act
	err := api.mapPin("TestBoard1", 1, "device 1",
		func(string) (freePins APIPinsMap) { return make(APIPinsMap) })
	// assert
	require.NotNil(err)
	assert.Contains(err.Error(), "Rail device")
}

func TestMapPinWithoutFreePinGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	bm := make(BoardsMap)
	bm["TestBoard1"] = &boards1Mock{name: "TestBoard1"}
	// mapped pins
	api := &BoardsAPI{
		mappedPins: make(APIPinsMap),
		boards:     bm,
	}
	// act
	err := api.mapPin("TestBoard1", 0, "device",
		func(string) (freePins APIPinsMap) { return make(APIPinsMap) })
	// assert
	require.NotNil(err)
	assert.Contains(err.Error(), "No free pin at")
}

func TestMapPinWithAlreadyMappedBoardPinGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	bm := make(BoardsMap)
	bm["TestBoard1"] = &boards1Mock{name: "TestBoard1"}
	// mapped pins
	mps := make(APIPinsMap)
	mps["k1"] = &apiPin{boardID: "TestBoard1", boardPinNr: 0, railDeviceName: "bin device 1"}
	api := &BoardsAPI{
		mappedPins: mps,
		boards:     bm,
	}
	// act
	err := api.mapPin("TestBoard1", 0, "bin device 2",
		func(string) (freePins APIPinsMap) {
			return APIPinsMap{"freekey": &apiPin{boardID: "TestBoard1", boardPinNr: 1}}
		})
	// assert
	require.NotNil(err)
	assert.Contains(err.Error(), "Pin already")
}

func TestReleasePin(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	// mapped pins
	mps := make(APIPinsMap)
	mps["device_1"] = &apiPin{railDeviceName: "device 1"}
	api := &BoardsAPI{
		mappedPins: mps,
		boards:     make(BoardsMap),
	}
	// act
	err := api.ReleasePin("device 1")
	// assert
	require.Nil(err)
	assert.Equal(0, len(api.mappedPins))
}

func TestReleasePinOfNotMappedRailDeviceGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	// mapped pins
	mps := make(APIPinsMap)
	mps["device_1"] = &apiPin{railDeviceName: "device 1"}
	api := &BoardsAPI{
		mappedPins: mps,
		boards:     make(BoardsMap),
	}
	// act
	err := api.ReleasePin("device 2")
	// assert
	require.NotNil(err)
	assert.Contains(err.Error(), "not mapped")
	assert.Contains(err.Error(), "device 2")
}

func TestReleasePinWithoutMappedPinsGetsError(t *testing.T) {
	// arrange
	assert := assert.New(t)
	require := require.New(t)
	api := &BoardsAPI{
		mappedPins: make(APIPinsMap),
		boards:     make(BoardsMap),
	}
	// act
	err := api.ReleasePin("a device")
	// assert
	require.NotNil(err)
	assert.Contains(err.Error(), "not mapped")
	assert.Contains(err.Error(), "a device")
}

func (b boards1Mock) GobotDevices() []gobot.Device { return nil }
func (b boards1Mock) GetBinaryPinNumbers() map[uint8]struct{} {
	return createBoardPinNumbersMap(0, b.binPins)
}
func (b boards1Mock) GetAnalogPinNumbers() map[uint8]struct{} {
	return createBoardPinNumbersMap(b.binPins, b.anaPins)
}
func (b boards1Mock) GetMemoryPinNumbers() map[uint8]struct{} {
	return createBoardPinNumbersMap(b.binPins+b.anaPins, b.memPins)
}

func createBoardPinNumbersMap(offset uint8, countEntries uint8) (fm map[uint8]struct{}) {
	fm = make(map[uint8]struct{})
	for i := offset; i < countEntries+offset; i++ {
		fm[i] = struct{}{}
	}
	return fm
}

// not used in this test scenarios
func (b boards1Mock) ReadValue(boardPinNr uint8) (uint8, error)          { return 0, nil }
func (b boards1Mock) SetValue(boardPinNr uint8, value uint8) (err error) { return }
func (b boards1Mock) ShowBoardConfig()                                   { return }
