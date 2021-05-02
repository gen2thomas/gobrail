package main

import (
	"context"
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

	if err = conf.Fill(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	if err = run(ctx, conf); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Println("end")
}

func reinit(c *config.Config) (err error) {
	rail, err = gobrailcreator.Create(true, "Model railroad prototype", c.AdaptorType, c.PlanFile, gobrailcreator.RecipeFiles{})
	return
}

func run(ctx context.Context, c *config.Config) (err error) {
	if err = reinit(c); err != nil {
		return
	}
	// https://forum.golangbridge.org/t/runtime-siftdowntimer-consuming-60-of-the-cpu/3773
	ticker := time.NewTicker(c.Tick)
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
