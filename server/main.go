package server

import (
	"crypto/sha512"
	"crypto/subtle"
	"path"
	"sync"
	"time"

	"github.com/SkyeYoung/url-screenshot-service/internal/helper"
	"github.com/SkyeYoung/url-screenshot-service/internal/r2"
	"github.com/SkyeYoung/url-screenshot-service/internal/screenshot"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/utils"
)

func Start(cfg *helper.Config, wg *sync.WaitGroup) {
	defer wg.Done()
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
	app.Use(cache.New(cache.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Query("refresh") == "true"
		},
		KeyGenerator: func(c *fiber.Ctx) string {
			return utils.CopyString(c.Path() + string(c.Request().Body()))
		},
		Expiration:   24 * time.Hour,
		CacheControl: true,
	}))

	// root routes
	if cfg.EnableMetrics {
		app.Get("/metrics", monitor.New())
	}

	// screenshot routes
	r2 := r2.New(cfg)
	logger := helper.GetLogger("server")
	screenshotApi := app.Group("/screenshot")
	screenshotPool := screenshot.Pool(cfg.Prefix)

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
		if res := screenshotPool.Process(url).(screenshot.Response); res.Err != nil {
			return res.Err
		}

		info, err := r2.UploadObject(&key)
		if err != nil {
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
