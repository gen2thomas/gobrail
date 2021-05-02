// Package config package manage program configuration by command line parameters
package config

import (
	"flag"
	"log"
	"time"

	"github.com/gen2thomas/gobrail/internal/app/gobrailcreator"
)

const defaultTick = 10 * time.Millisecond
const defaultPlan = "./plans/plan.json"

// Config contains program parameters
type Config struct {
	PlanFile    string
	AdaptorType gobrailcreator.AdaptorType
	Tick        time.Duration
}

// Fill will parse the command line and fill the configuration object
func (c *Config) Fill() (err error) {
	// this supports one plan (no further files will read)
	planFile := flag.String("plan", defaultPlan, "Path to railroad plan file")
	adaptorType := flag.String("adaptor", defaultAdaptor, "Supported adaptors are "+supportedAdaptors)
	tick := flag.Duration("tick", defaultTick, "Ticking interval, 10ms ... 50ms would be sufficient")
	flag.Parse()

	c.PlanFile = *planFile
	log.Println(c.PlanFile)
	c.AdaptorType, err = gobrailcreator.ParseAdaptorType(*adaptorType)
	c.Tick = *tick
	return
}
