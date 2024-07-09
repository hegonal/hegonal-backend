package controllers

import (
	"database/sql"

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

func NewTeamsInvite(c *fiber.Ctx) error {
	createNewTeamInvite := &models.CreateNewTeamInvite{}

	userID := c.Cookies("userID")

	if err := c.BodyParser(createNewTeamInvite); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	validate := utils.NewValidator()
	if err := validate.Struct(createNewTeamInvite); err != nil {
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

	teamMember, err := db.GetTeamMember(userID, createNewTeamInvite.TeamID)

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

	if teamMember.Role < models.TeamAdmin || createNewTeamInvite.Role >= teamMember.Role || createNewTeamInvite.Role == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Unauthorized",
		})
	}

	invitedUser, err := db.GetUserByEmail(createNewTeamInvite.Email)
	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Can't find this user",
		})
	} else if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	_, err = db.GetTeamMember(invitedUser.UserID, createNewTeamInvite.TeamID)
	if err != nil && err != sql.ErrNoRows {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	if err != sql.ErrNoRows {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "This user already in your team",
		})
	}

	_, err = db.GetTeamInviteByUserIDAndTeamID(createNewTeamInvite.TeamID, invitedUser.UserID)
	if err != nil && err != sql.ErrNoRows {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	if err != sql.ErrNoRows {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "You already send invite to this user.",
		})
	}

	teamInvite := models.TeamInvite{}
	teamInvite.InviteID = utils.GenerateId()
	teamInvite.TeamID = createNewTeamInvite.TeamID
	teamInvite.UserID = invitedUser.UserID
	teamInvite.Role = createNewTeamInvite.Role
	teamInvite.ExpiryDate = createNewTeamInvite.ExpiryDate
	teamInvite.CreatedAt = utils.TimeNow()
	teamInvite.UpdatedAt = utils.TimeNow()

	if err := db.CreateNewTeamInvite(&teamInvite); err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":       false,
		"msg":         nil,
		"team_invite": teamInvite,
	})
}

func AccpetInvite(c *fiber.Ctx) error {
	invitedID := c.Params("inviteID")

	userID := c.Cookies("userID")

	db, ok := c.Locals("db").(*database.Queries)
	if !ok {
		log.Error("Failed to retrieve DB from context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Failed to retrieve DB from context",
		})
	}

	teamInvite, err := db.GetTeamInvite(userID, invitedID)

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

	if teamInvite.ExpiryDate != nil && teamInvite.ExpiryDate.Before(utils.TimeNow()) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "This invite is already expired",
		})
	}

	teamMember := models.TeamMember{}
	teamMember.MemberID = userID
	teamMember.TeamID = teamInvite.TeamID
	teamMember.Role = teamInvite.Role
	teamMember.CreatedAt = utils.TimeNow()
	teamMember.UpdatedAt = utils.TimeNow()
	if err := db.CreateNewMember(&teamMember); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})

	}

	if err := db.DeleteTeamInviteByID(invitedID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":       false,
		"msg":         nil,
		"team_member": teamMember,
	})
}

func GetAllTeamMembers(c *fiber.Ctx) error {
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

	teamMembers, err := db.GetAllTeamMembersWithDetails(teamID)
	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	if len(teamMembers) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "This team is not exist.",
		})
	}

	userIsTeamMember := false

	for _, teamMember := range teamMembers {
		if teamMember.MemberID == userID {
			userIsTeamMember = true
		}
	}

	if !userIsTeamMember {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Unauthorized",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":        false,
		"msg":          nil,
		"team_members": teamMembers,
	})
}
func GetAllTeams(c *fiber.Ctx) error {
	userID := c.Cookies("userID")

	db, ok := c.Locals("db").(*database.Queries)
	if !ok {
		log.Error("Failed to retrieve DB from context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Failed to retrieve DB from context",
		})
	}

	teamMembers, err := db.GetUserAllTeams(userID)
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
		"user_teams": teamMembers,
	})
}
