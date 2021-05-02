package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gen2thomas/gobrail/internal/app/config"
	"github.com/gen2thomas/gobrail/internal/app/gobrailcreator"
)

var rail gobrailcreator.RailRunner

func main() {
	conf := &config.Config{}
	log.SetOutput(os.Stdout)
	var err error

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	reloadChan := make(chan bool, 1)

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
					reloadChan <- true
				}
			case <-ctx.Done():
				log.Printf("Done.")
				os.Exit(0)
			}
		}
	}()

	if err = conf.Fill(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	if err = reinit(*conf); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	if err = run(ctx, *conf, reloadChan); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Println("end")
}

func reinit(c config.Config) (err error) {
	rail, err = gobrailcreator.Create(true, "Model railroad prototype", c.AdaptorType, c.PlanFile, gobrailcreator.RecipeFiles{})
	return
}

func run(ctx context.Context, conf config.Config, reloadChan <-chan bool) (err error) {
	// https://forum.golangbridge.org/t/runtime-siftdowntimer-consuming-60-of-the-cpu/3773
	ticker := time.NewTicker(conf.Tick)
	// for monitor the cycle time
	firstRun := true
	var start time.Time
	// first loops takes >1 sec, so skip it for measurement
	count := -1
	const testCycles = 100
	for {
		select {
		case <-ctx.Done():
			fmt.Println("done in run")
			return nil
		case <-reloadChan:
			ticker.Stop()
			log.Printf("Wait finishing loop")
			if err = reinit(conf); err != nil {
				fmt.Errorf("Create rail with new configuration has error: %s", err.Error())
				return
			}
			count = -1
			ticker = time.NewTicker(conf.Tick)
		case <-ticker.C:
			count++
			if count == 1 {
				start = time.Now()
			}
			if err = rail.Run(); err != nil {
				return
			}
			if count == testCycles {
				count = -1
				ct := time.Since(start) / testCycles
				if firstRun {
					firstRun = false
					fmt.Printf("Cyclic time after %d runs with tick %s: %s\n", testCycles, conf.Tick, ct)
				}
				if ct > conf.Tick/10*18 {
					fmt.Printf("Warning: Cyclic time within last %d runs with tick %s: %s\n", testCycles, conf.Tick, ct)
				}
			}
		}
	}
}
