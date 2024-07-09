package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hegonal/hegonal-backend/app/controllers"
	"github.com/hegonal/hegonal-backend/pkg/middleware"
)

func TeamRoutes(a fiber.Router) {
	teamGroup := a.Group("/team")
	teamGroup.Use(middleware.SessionValidationMiddleware)

	teamGroup.Post("/add", controllers.CreateNewTeam)
	teamGroup.Post("/invite", controllers.NewTeamsInvite)

	teamGroup.Put("/invite/accpet/:inviteID", controllers.AccpetInvite)
	
	teamGroup.Get("/members/:teamID", controllers.GetAllTeamMembers)
	teamGroup.Get("/teams", controllers.GetAllTeams)
}
