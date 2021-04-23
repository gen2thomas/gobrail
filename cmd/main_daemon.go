package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gen2thomas/gobrail/internal/app/gobrailcreator"
)

const defaultTick = 10 * time.Millisecond
const defaultPlan = "./plans/plan.json"

//const defaultAdaptor = "Raspi"
const defaultAdaptor = "Digispark"

type config struct {
	planFile    string
	adaptorType gobrailcreator.AdaptorType
	tick        time.Duration
}

var rail gobrailcreator.RailRunner

func main() {
	conf := &config{}
	log.SetOutput(os.Stdout)
	var err error

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	defer func() {
		signal.Stop(signalChan)
		cancel()
	}()

	go func() {
		for {
			select {
			case s := <-signalChan:
				switch s {
				case syscall.SIGINT, syscall.SIGTERM:
					log.Printf("Got SIGINT/SIGTERM, exiting.")
					if err = gobrailcreator.Stop(); err != nil {
						log.Printf("Error while stopping: %s", err.Error())
					}
					cancel()
					os.Exit(1)
				case syscall.SIGHUP:
					log.Printf("Got SIGHUP, reloading.")
					if err = reinit(conf); err != nil {
						log.Printf("Create rail with new configuration has error: %s", err.Error())
						os.Exit(1)
					}
				}
			case <-ctx.Done():
				log.Printf("Done.")
				os.Exit(0)
			}
		}
	}()

	if err = conf.fill(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	if err = run(ctx, conf); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Println("end")
}

func (c *config) fill() (err error) {
	// this supports one plan (no further files will read)
	planFile := flag.String("plan", defaultPlan, "Path to railroad plan file")
	adaptorType := flag.String("adaptor", defaultAdaptor, "Supported adaptors are 'Digispark', 'Raspi', 'Tinkerboard'")
	tick := flag.Duration("tick", defaultTick, "Ticking interval")
	flag.Parse()

	c.planFile = *planFile
	log.Println(c.planFile)
	c.adaptorType, err = gobrailcreator.ParseAdaptorType(*adaptorType)
	c.tick = *tick
	return
}

func reinit(c *config) (err error) {
	rail, err = gobrailcreator.Create(true, "Model railroad prototype", c.adaptorType, c.planFile, gobrailcreator.RecipeFiles{})
	return
}

func run(ctx context.Context, c *config) (err error) {
	if err = reinit(c); err != nil {
		return
	}
	// https://forum.golangbridge.org/t/runtime-siftdowntimer-consuming-60-of-the-cpu/3773
	ticker := time.NewTicker(c.tick)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err = rail.Run(); err != nil {
				return
			}
		}
	}
}
