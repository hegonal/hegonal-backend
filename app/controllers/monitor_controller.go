package controllers

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/hegonal/hegonal-backend/app/models"
	"github.com/hegonal/hegonal-backend/pkg/utils"
	"github.com/hegonal/hegonal-backend/platform/database"
)

func CreateNewHttpMonitor(c *fiber.Ctx) error {
	createNewHttpMonitor := &models.CreateNewHttpMonitor{}

	userID := c.Cookies("userID")

	if err := c.BodyParser(createNewHttpMonitor); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	validate := utils.NewValidator()
	if err := validate.Struct(createNewHttpMonitor); err != nil {
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

	teamMember, err := db.GetTeamMember(userID, createNewHttpMonitor.TeamID)

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

	httpMonitor := &models.HttpMonitor{}
	httpMonitor.HttpMonitorID = utils.GenerateId()
	httpMonitor.TeamID = createNewHttpMonitor.TeamID
	httpMonitor.Status = models.HttpMonitorStatusUnknow
	httpMonitor.URL = createNewHttpMonitor.URL
	httpMonitor.Interval = createNewHttpMonitor.Interval
	httpMonitor.Retries = createNewHttpMonitor.Retries
	httpMonitor.RetryInterval = createNewHttpMonitor.RetryInterval
	httpMonitor.RequestTimeout = createNewHttpMonitor.RequestTimeout
	httpMonitor.ResendNotification = createNewHttpMonitor.ResendNotification
	httpMonitor.FollowRedirections = createNewHttpMonitor.FollowRedirections
	httpMonitor.MaxRedirects = createNewHttpMonitor.MaxRedirects
	httpMonitor.CheckSslError = createNewHttpMonitor.CheckSslError
	httpMonitor.SslExpiryReminders = createNewHttpMonitor.SslExpiryReminders
	httpMonitor.DomainExpiryReminders = createNewHttpMonitor.DomainExpiryReminders
	httpMonitor.HttpStatusCodes = createNewHttpMonitor.HttpStatusCodes
	httpMonitor.HttpMethod = createNewHttpMonitor.HttpMethod
	httpMonitor.BodyEncoding = createNewHttpMonitor.BodyEncoding
	httpMonitor.RequestBody = createNewHttpMonitor.RequestBody
	httpMonitor.RequestHeaders = createNewHttpMonitor.RequestHeaders
	httpMonitor.Group = createNewHttpMonitor.Group
	httpMonitor.Proxy = createNewHttpMonitor.Proxy
	httpMonitor.SendToOnCall = createNewHttpMonitor.SendToOnCall
	httpMonitor.CreatedAt = utils.TimeNow()
	httpMonitor.UpdatedAt = utils.TimeNow()

	if err := validate.Struct(createNewHttpMonitor); err != nil {
		log.Error(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	if err := db.CreateNewHttpMonitor(httpMonitor); err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":        false,
		"msg":          nil,
		"http_monitor": httpMonitor,
	})
}

func GetAllHttpMonitors(c *fiber.Ctx) error {
	teamID := c.Params("teamID")

	userID := c.Cookies("userID")

	db, ok := c.Locals("db").(*database.Queries)
	if !ok {
		log.Error("Failed to retrieve DB from context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Failed to retrieve DB from context",
		})
	}

	_, err := db.GetTeamMember(userID, teamID)
	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Unauthorized",
		})
	}else if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	httpMonitors, err := db.GetAllTeamHttpMonitors(teamID)

	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":        false,
		"msg":          nil,
		"http_monitors": httpMonitors,
	})
}

func DeleteHttpMonitor(c *fiber.Ctx) error {
	teamID := c.Params("teamID")
	monitorID := c.Params("monitorID")

	userID := c.Cookies("userID")

	db, ok := c.Locals("db").(*database.Queries)
	if !ok {
		log.Error("Failed to retrieve DB from context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Failed to retrieve DB from context",
		})
	}

	teamMember, err := db.GetTeamMember(userID, teamID)
	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Unauthorized",
		})
	}else if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	if teamMember.Role < models.TeamAdmin {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Unauthorized",
		})
	}

	err = db.DeleteHttpMonitor(monitorID)

	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}