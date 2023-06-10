package server

import (
	"errors"

	"github.com/SkyeYoung/url-screenshot-service/internal/helper"
	"github.com/gofiber/fiber/v2"
)

func getUrlFromRequest(c *fiber.Ctx) (string, error) {
	d := new(Website)
	if err := c.BodyParser(d); err != nil {
		return "", errors.New("invalid url")
	}

	return helper.GetValidUrl(d.Url)
}
