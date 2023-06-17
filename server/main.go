package server

import (
	"path"

	"github.com/SkyeYoung/url-screenshot-service/internal/helper"
	"github.com/SkyeYoung/url-screenshot-service/internal/r2"
	"github.com/SkyeYoung/url-screenshot-service/internal/screenshot"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func Start(cfg *helper.Config) {
	app := fiber.New()

	// middleware
	setupMiddleware(app, cfg)

	// root routes
	if cfg.EnableMetrics {
		app.Get("/metrics", monitor.New())
	}

	// screenshot routes
	r2 := r2.New(cfg)
	logger := helper.GetLogger("server")
	screenshotApi := app.Group("/screenshot")
	ss := screenshot.New(cfg.Prefix)
	defer ss.Close()

	screenshotApi.Post("/", func(c *fiber.Ctx) error {
		url, err := getUrlFromRequest(c)
		if err != nil {
			logger.Error(err)
			return err
		}
		logger.Infof("processing request from %v for %v", c.IP(), url)

		key := path.Join(cfg.Prefix, helper.EncodeImgNameAddExt(url))
		logger.Infof("%v generated key: %v", url, key)

		logger.Info("checking if screenshot key exists")
		if _, err := r2.HeadObject(&key); err == nil {
			logger.Infof("screenshot already exists, returning url: `%v`", url)
			return c.SendString(cfg.ReturnUrl + "/" + key)
		} else {
			logger.Infof("screenshot does not exist, because: %v", err)
		}

		logger.Infof("trying to get screeshot of %v", url)
		if res := ss.GetPool().Process(url); res != nil {
			err := res.(*screenshot.Response).Err
			logger.Error(err)
			return err
		}

		info, err := r2.UploadObject(&key)
		if err != nil {
			logger.Error(err)
			return err
		}
		logger.Infof("screenshot uploaded to %v", info.Location)

		if err := helper.RmImgAfterUpload(cfg, logger, url, key); err != nil {
			return err
		}

		logger.Infof("returning url: `%v`", url)
		return c.SendString(cfg.ReturnUrl + "/" + key)
	})

	screenshotApi.Delete("/", func(c *fiber.Ctx) error {
		url, err := getUrlFromRequest(c)
		if err != nil {
			logger.Error(err)
			return err
		}
		logger.Infof("processing request from %v for %v", c.IP(), url)

		key := path.Join(cfg.Prefix, helper.EncodeImgNameAddExt(url))
		logger.Infof("%v generated key: %v", url, key)

		logger.Info("trying to delete screenshot key")
		_, err = r2.DeleteObject(&key)
		return err
	})

	logger.Fatal(app.Listen(":" + cfg.Port))
}
