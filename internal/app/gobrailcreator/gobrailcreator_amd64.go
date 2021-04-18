package gobrailcreator

// special implementation part for amd64 (x86) targets

import (
	"fmt"

	"gobot.io/x/gobot/platforms/digispark"
)

func createAdaptor(adaptorType AdaptorType) (adaptor i2cAdaptor, err error) {
	switch adaptorType {
	case digisparkType:
		adaptor = digispark.NewAdaptor()
	case raspiType:
		err = fmt.Errorf("Amd environment not supported by Raspi")
	case tinkerboardType:
		err = fmt.Errorf("Arm environment not supported by Tinkerboard")
	default:
		err = fmt.Errorf("Unknown type '%d'", adaptorType)
	}
	return
}
