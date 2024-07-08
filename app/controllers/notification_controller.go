package controllers

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/hegonal/hegonal-backend/app/models"
	"github.com/hegonal/hegonal-backend/pkg/utils"
	"github.com/hegonal/hegonal-backend/platform/database"
)

func CreateNewNotification(c *fiber.Ctx) error {
	createNotification := &models.CreateNotification{}

	userID := c.Cookies("userID")

	if err := c.BodyParser(createNotification); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	validate := utils.NewValidator()
	if err := validate.Struct(createNotification); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	if createNotification.TeamID == nil && createNotification.UserID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "You must specify at least one affiliation where you will use this notification.",
		})
	}

	if createNotification.TeamID != nil && createNotification.UserID != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "You cannot specify two institutions at the same time.",
		})
	}

	if createNotification.UserID != nil && createNotification.UserID != &userID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Unauthorized",
		})
	}

	db, ok := c.Locals("db").(*database.Queries)
	if !ok {
		log.Error("Failed to retrieve DB from context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Failed to retrieve DB from context",
		})
	}

	if createNotification.TeamID != nil {
		teamMember, err := db.GetTeamMember(userID, *createNotification.TeamID)

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

		if teamMember.Role < models.TeamAdmin {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": true,
				"msg":   "Unauthorized",
			})
		}
	}

	notification := &models.Notification{}
	notification.NotificationID = utils.GenerateId()
	notification.UserID = createNotification.UserID
	notification.TeamID = createNotification.TeamID
	notification.NotificationType = createNotification.NotificationType
	notification.NotificationConfig = createNotification.NotificationConfig
	notification.CreatedAt = utils.TimeNow()
	notification.UpdatedAt = utils.TimeNow()

	if err := db.CreateNewNotification(notification); err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":        false,
		"msg":          nil,
		"notification": notification,
	})
}

func CreateNewHttpMonitorNotification(c *fiber.Ctx) error {
	createHttpMonitorNotification := &models.CreateHttpMonitorNotification{}

	userID := c.Cookies("userID")

	if err := c.BodyParser(createHttpMonitorNotification); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	validate := utils.NewValidator()
	if err := validate.Struct(createHttpMonitorNotification); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	db, ok := c.Locals("db").(*database.Queries)
	if !ok {
		log.Error("Failed to retrieve DB from context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Failed to retrieve DB from context",
		})
	}

	teamMember, err := db.GetTeamMember(userID, createHttpMonitorNotification.TeamID)

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

	if teamMember.Role < models.TeamAdmin {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Unauthorized",
		})
	}

	httpMonitorNotification := models.HttpMonitorNotification{}
	httpMonitorNotification.HttpMonitorID = createHttpMonitorNotification.HttpMonitorID
	httpMonitorNotification.NotificationID = createHttpMonitorNotification.NotificationID

	if err := db.CreateNewHttpMonitorNotification(&httpMonitorNotification); err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":        false,
		"msg":          nil,
	})
}
