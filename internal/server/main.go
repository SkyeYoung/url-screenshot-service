package server

import (
	"crypto/sha512"
	"crypto/subtle"
	"path"

	"github.com/SkyeYoung/url-screenshot-service/internal/helper"
	"github.com/SkyeYoung/url-screenshot-service/internal/r2"
	"github.com/SkyeYoung/url-screenshot-service/internal/screenshot"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func Start(cfg *helper.Config) {
	app := fiber.New()

	// middleware
	app.Use(helmet.New())
	app.Use(idempotency.New())
	app.Use(keyauth.New(keyauth.Config{
		Validator: func(c *fiber.Ctx, key string) (bool, error) {
			hashedKey := sha512.Sum512([]byte(key))
			hashedApiKey := sha512.Sum512([]byte(cfg.ApiKey))

			if subtle.ConstantTimeCompare(hashedKey[:], hashedApiKey[:]) == 1 {
				return true, nil
			}
			return false, keyauth.ErrMissingOrMalformedAPIKey
		},
	}))

	// routes
	if cfg.EnableMetrics {
		app.Get("/metrics", monitor.New())
	}

	r2 := r2.New(cfg)
	logger := helper.GetLogger("server")
	app.Post("/", func(c *fiber.Ctx) error {
		url, err := getUrlFromRequest(c)
		if err != nil {
			logger.Error(err)
			return err
		}
		logger.Infof("processing request from %v for %v", c.IP(), url)

		key := path.Join(cfg.Prefix, helper.EncodeImgNameAddExt(url))
		logger.Info(url + " generated key: " + key)

		logger.Info("checking if screenshot key exists")
		if _, err := r2.GetObject(&key); err == nil {
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
		if info, err := r2.UploadObject(&key); err != nil {
			return err
		} else {
			logger.Infof("screenshot uploaded to %v", info.Location)
			return c.SendString(url)
		}
	})

	logger.Fatal(app.Listen(":" + cfg.Port))
}
