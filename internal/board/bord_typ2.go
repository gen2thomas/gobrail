package board

/* Implementation for circuit board "Typ2" with one I2C chip PCA9501
 *
 *      Author: g2t
 *  Created on: 01.06.2009 (in C++)
 *    Modified: 19.04.2013
 *   in golang: 28.03.2021
 * Called from: Modellbahn.cpp (outdated)
 * Call       : some functions from i2c (PCA9501), digispark
 *
 * 9501:
 * - 8 GPIO, 0..3 amplified and negotiated with IRLZ34N, 4..7 for max. 20mA (see docs/images)
 * - Dummy address for EEPROM write is set to 0x00
 *
 * Functions:
 * + read eeprom at board
 * + read/write eeprom at board for "configmode"
 *
 * TODO:
 * - provide EEPROM as separat "chip" of this board
 * - ensure no overlap of EEPROM for general use with "configmode"
 */

import (
	"fmt"
	"time"

	"gobot.io/x/gobot/drivers/i2c"
)

const boardTyp2IoNrMax = 7 //0..7
const chipsAtBoardTyp2Max = 1
const lastEepromAddressForConfigmode = uint8(0xFF)
const chipID = "PCA9501.GPIO.Mem"

//this is the default io configuration of this board
var boardPinsDefault = PinsMap{
	0:  {chipID: chipID, chipPin: 0, pinType: Binary},
	1:  {chipID: chipID, chipPin: 1, pinType: Binary},
	2:  {chipID: chipID, chipPin: 2, pinType: Binary},
	3:  {chipID: chipID, chipPin: 3, pinType: Binary},
	4:  {chipID: chipID, chipPin: 4, pinType: Binary},
	5:  {chipID: chipID, chipPin: 5, pinType: Binary},
	6:  {chipID: chipID, chipPin: 6, pinType: Binary},
	7:  {chipID: chipID, chipPin: 7, pinType: Binary},
	8:  {chipID: chipID, chipPin: 0x01, pinType: Memory},
	9:  {chipID: chipID, chipPin: 0x02, pinType: Memory},
	10: {chipID: chipID, chipPin: 0x02, pinType: Memory},
	11: {chipID: chipID, chipPin: 0x03, pinType: Memory},
	12: {chipID: chipID, chipPin: 0x04, pinType: Memory},
	13: {chipID: chipID, chipPin: 0x05, pinType: Memory},
	14: {chipID: chipID, chipPin: 0x06, pinType: Memory},
	15: {chipID: chipID, chipPin: 0x07, pinType: Memory},
}

// NewBoardTyp2 creates a new board of typ 2
func NewBoardTyp2(adaptor i2c.Connector, address uint8, name string) *Board {
	p := &Board{
		name: name,
		pins: boardPinsDefault,
		chips: map[string]*chip{chipID: {
			address: address,
			device:  i2c.NewPCA9501Driver(adaptor, i2c.WithAddress(int(address))),
		}},
	}

	return p
}

// WriteBoardConfig is writing the configuration to EEPROM
func (b *Board) WriteBoardConfig() error {
	eeaddress := lastEepromAddressForConfigmode
	// this will only work if wc pin is high!
	// write the IO's
	for ioNr, boardPin := range b.pins {
		err := b.writeEEPROM(eeaddress, ioNr)
		if err != nil {
			return err
		}
		eeaddress--
		err = b.writeEEPROM(eeaddress, boardPin.chipPin)
		if err != nil {
			return err
		}
		eeaddress--
	}

	return nil
}

// ReadBoardConfig is reading the configuration from EEPROM
func (b *Board) ReadBoardConfig() (err error) {
	eeaddress := lastEepromAddressForConfigmode
	// read the IO's
	for i := uint8(0); i < (boardTyp2IoNrMax + 1); i++ {
		ioNr, err := b.readEEPROM(eeaddress)
		if err != nil {
			return err
		}
		eeaddress--
		b.pins[ioNr].chipPin, err = b.readEEPROM(eeaddress)
		if err != nil {
			return err
		}
		eeaddress--
	}
	return nil
}

func (b *Board) writeGPIO(boardPinNr uint8, val uint8) (err error) {
	var pin *boardPin
	var ok bool
	if pin, ok = b.pins[boardPinNr]; !ok {
		err = fmt.Errorf("There is no pin with key '%d' for writeGPIO", boardPinNr)
		return
	}
	var params = map[string]interface{}{
		"pin": pin.chipPin,
		"val": val,
	}
	writeGpioCommand := b.chips[pin.chipID].device.Command("WriteGPIO")
	result := writeGpioCommand(params).(map[string]interface{})["err"]
	if result != nil {
		return result.(error)
	}
	return
}

func (b *Board) readGPIO(boardPinNr uint8) (val uint8, err error) {
	var pin *boardPin
	var ok bool
	if pin, ok = b.pins[boardPinNr]; !ok {
		err = fmt.Errorf("There is no pin with key '%d' for readGPIO", boardPinNr)
		return
	}
	params := make(map[string]interface{})
	params["pin"] = pin.chipPin
	readGpioCommand := b.chips[pin.chipID].device.Command("ReadGPIO")
	result := readGpioCommand(params).(map[string]interface{})
	if result["err"] != nil {
		return 0, result["err"].(error)
	}
	return result["val"].(uint8), nil
}

func (b *Board) writeEEPROM(address uint8, val uint8) (err error) {
	var pin *boardPin
	var ok bool
	if pin, ok = b.pins[address]; !ok {
		err = fmt.Errorf("There is no pin with key '%d' for writeEEPROM", address)
		return
	}
	var params = map[string]interface{}{
		"address": pin.chipPin,
		"val":     val,
	}
	writeMemCommand := b.chips[pin.chipID].device.Command("WriteEEPROM")
	result := writeMemCommand(params).(map[string]interface{})["err"]
	time.Sleep(4 * time.Millisecond)
	if result != nil {
		return result.(error)
	}
	return
}

func (b *Board) readEEPROM(address uint8) (val uint8, err error) {
	var pin *boardPin
	var ok bool
	if pin, ok = b.pins[address]; !ok {
		err = fmt.Errorf("There is no pin with key '%d' for readEEPROM", address)
		return
	}
	params := make(map[string]interface{})
	params["address"] = pin.chipPin
	readMemCommand := b.chips[pin.chipID].device.Command("ReadEEPROM")
	result := readMemCommand(params).(map[string]interface{})
	time.Sleep(4 * time.Millisecond)
	if result["err"] != nil {
		return 0, result["err"].(error)
	}
	return result["val"].(uint8), nil
}
