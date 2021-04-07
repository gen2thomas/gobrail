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
// + read/write eeprom at board
// + read/write GPIO at board
//
// TODO:
// - provide EEPROM as separat "chip" of this board
// - config with "from - to" range, especially for memory
// - ensure no overlap of EEPROM for general use with "configmode"
//

import (
	"time"

	"gobot.io/x/gobot/drivers/i2c"
)

const chipID = "PCA9501.GPIO.Mem"

//this is the default io configuration of this board
var boardPinsDefault = PinsMap{
	0:  {chipID: chipID, chipPinNr: 0, pinType: BinaryW},
	1:  {chipID: chipID, chipPinNr: 1, pinType: BinaryW},
	2:  {chipID: chipID, chipPinNr: 2, pinType: BinaryW},
	3:  {chipID: chipID, chipPinNr: 3, pinType: BinaryW},
	4:  {chipID: chipID, chipPinNr: 4, pinType: NBinaryR},
	5:  {chipID: chipID, chipPinNr: 5, pinType: NBinaryR},
	6:  {chipID: chipID, chipPinNr: 6, pinType: NBinaryR},
	7:  {chipID: chipID, chipPinNr: 7, pinType: NBinaryR},
	8:  {chipID: chipID, chipPinNr: 0x01, pinType: Memory},
	9:  {chipID: chipID, chipPinNr: 0x02, pinType: Memory},
	10: {chipID: chipID, chipPinNr: 0x02, pinType: Memory},
	11: {chipID: chipID, chipPinNr: 0x03, pinType: Memory},
	12: {chipID: chipID, chipPinNr: 0x04, pinType: Memory},
	13: {chipID: chipID, chipPinNr: 0x05, pinType: Memory},
	14: {chipID: chipID, chipPinNr: 0x06, pinType: Memory},
	15: {chipID: chipID, chipPinNr: 0x07, pinType: Memory},
}

// NewBoardTyp2 creates a new board of typ 2
func NewBoardTyp2(adaptor i2c.Connector, address uint8, name string) *Board {
	chips := map[string]*chip{chipID: {
		address: address,
		driver:  i2c.NewPCA9501Driver(adaptor, i2c.WithAddress(int(address))),
	}}

	return NewBoard(name, chips, boardPinsDefault)
}

func (b *Board) writeGPIO(bPin *boardPin, val uint8) (err error) {
	var driver DriverOperations
	if driver, err = b.getDriver(bPin); err != nil {
		return
	}
	var params = map[string]interface{}{
		"pin": bPin.chipPinNr,
		"val": val,
	}
	result := driver.Command("WriteGPIO")(params).(map[string]interface{})["err"]
	if result != nil {
		return result.(error)
	}
	return
}

func (b *Board) readGPIO(bPin *boardPin) (val uint8, err error) {
	var driver DriverOperations
	if driver, err = b.getDriver(bPin); err != nil {
		return
	}
	params := make(map[string]interface{})
	params["pin"] = bPin.chipPinNr
	result := driver.Command("ReadGPIO")(params).(map[string]interface{})
	if result["err"] != nil {
		return 0, result["err"].(error)
	}
	return result["val"].(uint8), nil
}

func (b *Board) writeEEPROM(bPin *boardPin, val uint8) (err error) {
	var driver DriverOperations
	if driver, err = b.getDriver(bPin); err != nil {
		return
	}
	var params = map[string]interface{}{
		"address": bPin.chipPinNr,
		"val":     val,
	}
	result := driver.Command("WriteEEPROM")(params).(map[string]interface{})["err"]
	time.Sleep(4 * time.Millisecond)
	if result != nil {
		return result.(error)
	}
	return
}

func (b *Board) readEEPROM(bPin *boardPin) (val uint8, err error) {
	var driver DriverOperations
	if driver, err = b.getDriver(bPin); err != nil {
		return
	}
	params := make(map[string]interface{})
	params["address"] = bPin.chipPinNr
	result := driver.Command("ReadEEPROM")(params).(map[string]interface{})
	time.Sleep(4 * time.Millisecond)
	if result["err"] != nil {
		return 0, result["err"].(error)
	}
	return result["val"].(uint8), nil
}
