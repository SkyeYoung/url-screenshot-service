package main

import (
	"errors"
	"path"

	"github.com/SkyeYoung/url-screenshot-service/internal/helper"
	"github.com/SkyeYoung/url-screenshot-service/internal/r2"
	"github.com/SkyeYoung/url-screenshot-service/internal/screenshot"
	"github.com/gofiber/fiber/v2"
)

type Website struct {
	Url string `json:"url"`
}

func getUrlFromRequest(c *fiber.Ctx) (string, error) {
	d := new(Website)
	if err := c.BodyParser(d); err != nil {
		return "", errors.New("invalid url")
	}

	return helper.GetValidUrl(d.Url)
}

func main() {
	helper.SetupLogger()

	logger := helper.GetLogger()
	cfg, err := helper.GetConfig()
	r2.SetupSession(cfg)

	if err != nil {
		logger.Fatal(err)
	}

	app := fiber.New()
	app.Post("/", func(c *fiber.Ctx) error {
		url, err := getUrlFromRequest(c)
		if err != nil {
			logger.Error(err)
			return err
		}

		key := path.Join(cfg.Prefix, helper.EncodeImgNameAddExt(url))
		logger.Info(url + " generated key: " + key)
		logger.Info("checking if screenshot exists")
		if _, err := r2.GetObjectAttributes(cfg, &key); err == nil {
			logger.Info("screenshot already exists, returning url: " + url)
			return c.SendString(url)
		} else {
			logger.Info("screenshot does not exist, continuing")
			logger.Info(err)
		}

		logger.Infof("trying to get screeshot of %v", url)
		if _, err := screenshot.Screenshot(url, cfg.Prefix); err != nil {
			return err
		}
		if info, err := r2.Upload(cfg, &key); err != nil {
			return err
		} else {
			logger.Infof("screenshot uploaded to %v", info.Location)
			return c.SendString(url)
		}
	})

	logger.Fatal(app.Listen(":" + cfg.Port))
}
