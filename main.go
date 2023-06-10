package main

import (
	"github.com/SkyeYoung/url-screenshot-service/internal/helper"
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

	helper.GetLogger("server").Info("Starting server...")
	server.Start(cfg)
}
