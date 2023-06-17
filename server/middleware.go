package server

import (
	"crypto/sha512"
	"crypto/subtle"
	"time"

	"github.com/SkyeYoung/url-screenshot-service/internal/helper"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
	"github.com/gofiber/fiber/v2/utils"
)

func setupMiddleware(app *fiber.App, cfg *helper.Config) {
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
}
