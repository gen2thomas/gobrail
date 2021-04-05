package board

// NewBoardForTestWithoutChips creates a new board only for tests without any chip
func NewBoardForTestWithoutChips(name string, pinCountBin uint8, pinCountAna uint8, pinCountMem uint8) *Board {
	pins := make(PinsMap)

	for i := uint8(0); i < pinCountBin; i++ {
		pins[i] = &boardPin{chipID: "NoChip", chipPinNr: i, pinType: Binary}
	}
	for i := uint8(0); i < pinCountAna; i++ {
		id := pinCountBin + i
		pins[id] = &boardPin{chipID: "NoChip", chipPinNr: id, pinType: Analog}
	}
	for i := uint8(0); i < pinCountMem; i++ {
		id := pinCountBin + pinCountAna + i
		pins[id] = &boardPin{chipID: "NoChip", chipPinNr: id, pinType: Memory}
	}
	return NewBoard(name, nil, pins)
}
