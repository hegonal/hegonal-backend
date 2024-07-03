package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/hegonal/hegonal-backend/app/models"
	"github.com/hegonal/hegonal-backend/pkg/utils"
	"github.com/hegonal/hegonal-backend/platform/database"
)

func TeamAdd(c *fiber.Ctx) error {
	teamAdd := &models.TeamAdd{}

	_ = c.Cookies("userID")

	if err := c.BodyParser(teamAdd); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	validate := utils.NewValidator()
	if err := validate.Struct(teamAdd); err != nil {
		log.Error(err)
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

	team := &models.Team{}
	team.ID = utils.GenerateId()
	team.Name = teamAdd.Name
	team.Description = teamAdd.Description
	team.CreatedAt = time.Now().UTC()
	team.UpdatedAt = time.Now().UTC()

	if err := db.CreateNewTeam(team); err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"team":  team,
	})
}
