package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hegonal/hegonal-backend/app/controllers"
	"github.com/hegonal/hegonal-backend/pkg/middleware"
)

func TeamRoutes(a fiber.Router) {
	authGroup := a.Group("/team")
	authGroup.Use(middleware.SessionValidationMiddleware)

	authGroup.Post("/add", controllers.CreateNewTeam)
}
