package utils

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func ClearCookies(c *fiber.Ctx, key ...string) {
	for i := range key {
		c.Cookie(&fiber.Cookie{
			Name:    key[i],
			Expires: TimeNow().Add(-time.Hour * 24),
			Value:   "",
		})
	}
}
