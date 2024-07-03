package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/hegonal/hegonal-backend/app/models"
	"github.com/hegonal/hegonal-backend/app/queries"
	"github.com/hegonal/hegonal-backend/pkg/utils"
	"github.com/hegonal/hegonal-backend/platform/database"
)

func TeamAdd(c *fiber.Ctx) error {
	teamAdd := &models.TeamAdd{}

	userSession := c.Cookies("session")
	userId := c.Cookies("userID")

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

	db, err := database.OpenDBConnection()
	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	test := time.Now()
	newSession, err := db.RotateSession(userId, userSession)
	log.Info(time.Since(test))
	if _, ok := err.(*queries.NoRowsAffectedError); ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Please relogin or try again later",
		})
	} else if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "session",
		Value:    newSession,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	})

	team := &models.Team{}
	team.ID = utils.GenerateId()
	team.Name = teamAdd.Name
	team.Description = teamAdd.Description
	team.CreatedAt = time.Now()
	team.UpdatedAt = time.Now()

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
