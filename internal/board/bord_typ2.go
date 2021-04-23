package board

// Implementation for circuit board "Typ2" with one I2C chip PCA9501
//
//      Author: g2t
//  Created on: 01.06.2009 (in C++)
//    Modified: 19.04.2013
//   in golang: 28.03.2021
// Called from: boardsapi
// Call       : some functions from gobot-i2c (PCA9501)
//
// 9501:
// - 8 GPIO, 0..3 amplified and negotiated with IRLZ34N, 4..7 for max. 20mA (see docs/images)
// - Dummy address for EEPROM write is set to 0x00
//
// Functions:
// + read/write EEPROM at board
// + read/write GPIO at board
//
// TODO:
// - provide EEPROM as separate "chip" of this board
// - config with "from - to" range, especially for memory
// - ensure no overlap of EEPROM for general use with "configmode"
//

import (
	"time"

	"gobot.io/x/gobot/drivers/i2c"

	"github.com/gen2thomas/gobrail/internal/boardpin"
)

const chipID = "PCA9501.GPIO.Mem"

//this is the default io configuration of this board
var boardPinsDefault = PinsMap{
	0:  {ChipID: chipID, ChipPinNr: 0, PinType: boardpin.BinaryW},
	1:  {ChipID: chipID, ChipPinNr: 1, PinType: boardpin.BinaryW},
	2:  {ChipID: chipID, ChipPinNr: 2, PinType: boardpin.BinaryW},
	3:  {ChipID: chipID, ChipPinNr: 3, PinType: boardpin.BinaryW},
	4:  {ChipID: chipID, ChipPinNr: 4, PinType: boardpin.NBinaryR},
	5:  {ChipID: chipID, ChipPinNr: 5, PinType: boardpin.NBinaryR},
	6:  {ChipID: chipID, ChipPinNr: 6, PinType: boardpin.NBinaryR},
	7:  {ChipID: chipID, ChipPinNr: 7, PinType: boardpin.NBinaryR},
	8:  {ChipID: chipID, ChipPinNr: 0x01, PinType: boardpin.Memory},
	9:  {ChipID: chipID, ChipPinNr: 0x02, PinType: boardpin.Memory},
	10: {ChipID: chipID, ChipPinNr: 0x02, PinType: boardpin.Memory},
	11: {ChipID: chipID, ChipPinNr: 0x03, PinType: boardpin.Memory},
	12: {ChipID: chipID, ChipPinNr: 0x04, PinType: boardpin.Memory},
	13: {ChipID: chipID, ChipPinNr: 0x05, PinType: boardpin.Memory},
	14: {ChipID: chipID, ChipPinNr: 0x06, PinType: boardpin.Memory},
	15: {ChipID: chipID, ChipPinNr: 0x07, PinType: boardpin.Memory},
}

// NewBoardTyp2 creates a new board of type 2
func NewBoardTyp2(adaptor i2c.Connector, address uint8, name string) *Board {
	chips := map[string]*chip{chipID: {
		address: address,
		driver:  i2c.NewPCA9501Driver(adaptor, i2c.WithAddress(int(address))),
	}}

	return NewBoard(name, chips, boardPinsDefault)
}

func (b *Board) writeGPIO(bPin *boardpin.Pin, val uint8) (err error) {
	var driver DriverOperations
	if driver, err = b.getDriver(bPin); err != nil {
		return
	}
	var params = map[string]interface{}{
		"pin": bPin.ChipPinNr,
		"val": val,
	}
	result := driver.Command("WriteGPIO")(params).(map[string]interface{})["err"]
	if result != nil {
		return result.(error)
	}
	return
}

func (b *Board) readGPIO(bPin *boardpin.Pin) (val uint8, err error) {
	var driver DriverOperations
	if driver, err = b.getDriver(bPin); err != nil {
		return
	}
	params := make(map[string]interface{})
	params["pin"] = bPin.ChipPinNr
	result := driver.Command("ReadGPIO")(params).(map[string]interface{})
	if result["err"] != nil {
		return 0, result["err"].(error)
	}
	return result["val"].(uint8), nil
}

func (b *Board) writeEEPROM(bPin *boardpin.Pin, val uint8) (err error) {
	var driver DriverOperations
	if driver, err = b.getDriver(bPin); err != nil {
		return
	}
	var params = map[string]interface{}{
		"address": bPin.ChipPinNr,
		"val":     val,
	}
	result := driver.Command("WriteEEPROM")(params).(map[string]interface{})["err"]
	time.Sleep(4 * time.Millisecond)
	if result != nil {
		return result.(error)
	}
	return
}

func (b *Board) readEEPROM(bPin *boardpin.Pin) (val uint8, err error) {
	var driver DriverOperations
	if driver, err = b.getDriver(bPin); err != nil {
		return
	}
	params := make(map[string]interface{})
	params["address"] = bPin.ChipPinNr
	result := driver.Command("ReadEEPROM")(params).(map[string]interface{})
	time.Sleep(4 * time.Millisecond)
	if result["err"] != nil {
		return 0, result["err"].(error)
	}
	return result["val"].(uint8), nil
}
