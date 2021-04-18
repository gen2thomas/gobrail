package gobrailcreator

// special implementation part for arm targets

import (
	"fmt"

	"gobot.io/x/gobot/platforms/raspi"
	"gobot.io/x/gobot/platforms/tinkerboard"
)

func createAdaptor(adaptorType AdaptorType) (adaptor i2cAdaptor, err error) {
	switch adaptorType {
	case digisparkType:
		err = fmt.Errorf("Arm environment not supported by Digispark")
	case raspiType:
		adaptor = raspi.NewAdaptor()
	case tinkerboardType:
		adaptor = tinkerboard.NewAdaptor()
	default:
		err = fmt.Errorf("Unknown type '%d'", adaptorType)
	}
	return
}
