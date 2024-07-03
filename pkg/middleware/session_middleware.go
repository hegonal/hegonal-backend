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

	if userSession.ExpiryTime.Before(time.Now().UTC()) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Session expired, please login again",
		})
	}

	refreshThreshold := 23 * time.Hour + 45 * time.Minute
	if userSession.ExpiryTime.Sub(time.Now().UTC()) < refreshThreshold {
		newSession, err := db.RotateSession(userID, session)
		if err != nil {
			log.Error("Failed to rotate session:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   "Failed to rotate session",
			})
		}

		c.Cookie(&fiber.Cookie{
			Name:     "session",
			Value:    newSession,
			Expires:  time.Now().UTC().Add(24 * time.Hour),
			HTTPOnly: true,
			Secure:   true,
			SameSite: "Lax",
		})
	}
	
	c.Locals("db", db)

	return c.Next()
}
