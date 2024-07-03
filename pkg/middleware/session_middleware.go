package middleware

import (
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/hegonal/hegonal-backend/platform/database"
)

func SessionValidationMiddleware(c *fiber.Ctx) error {
	userID := c.Cookies("userID")
	session := c.Cookies("session")

	if userID == "" || session == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Unauthorized",
		})
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err,
		})
	}

	userSession, err := db.GetSession(userID, session)
	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Unauthorized",
		})
	} else if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err,
		})
	}

	if userSession.ExpiryTime.Before(time.Now()) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Session expired, please login again",
		})
	}

	return c.Next()
}
