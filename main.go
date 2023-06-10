package main

import (
	"sync"

	"github.com/SkyeYoung/url-screenshot-service/internal/helper"
	"github.com/SkyeYoung/url-screenshot-service/scheduler"
	"github.com/SkyeYoung/url-screenshot-service/server"
)

func main() {
	cfg, err := helper.SetupConfig()
	if err != nil {
		panic(err)
	}
	err = helper.SetupLogger(cfg)
	if err != nil {
		panic(err)
	}

	wg := new(sync.WaitGroup)
	wg.Add(2)

	helper.GetLogger("scheduler").Info("Starting scheduler...")
	go scheduler.New(cfg).Start(wg)

	helper.GetLogger("server").Info("Starting server...")
	go server.Start(cfg, wg)

	wg.Wait()
}
