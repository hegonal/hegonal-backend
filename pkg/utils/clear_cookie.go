package utils

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func ClearCookies(c *fiber.Ctx, key ...string) {
	for i := range key {
		c.Cookie(&fiber.Cookie{
			Name:    key[i],
			Expires: time.Now().Add(-time.Hour * 24),
			Value:   "",
		})
	}
}