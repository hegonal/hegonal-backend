package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/hegonal/hegonal-backend/app/models"
	"github.com/hegonal/hegonal-backend/pkg/utils"
	"github.com/hegonal/hegonal-backend/platform/database"
)

func CreateNewTeam(c *fiber.Ctx) error {
	createNewTeam := &models.CreateNewTeam{}

	userID := c.Cookies("userID")

	if err := c.BodyParser(createNewTeam); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	validate := utils.NewValidator()
	if err := validate.Struct(createNewTeam); err != nil {
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
	team.TeamID = utils.GenerateId()
	team.Name = createNewTeam.Name
	team.Description = createNewTeam.Description
	team.CreatedAt = utils.TimeNow()
	team.UpdatedAt = utils.TimeNow()

	if err := db.CreateNewTeam(team); err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	teamMember := &models.TeamMember{}
	teamMember.TeamID = team.TeamID
	teamMember.MemberID = userID
	teamMember.Role = models.TeamOwner
	teamMember.CreatedAt = utils.TimeNow()
	teamMember.UpdatedAt = utils.TimeNow()

	if err := db.CreateNewMember(teamMember); err != nil {
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
