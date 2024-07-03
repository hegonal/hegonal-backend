package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hegonal/hegonal-backend/app/controllers"
)

func TeamRoutes(a fiber.Router) {
	authGroup := a.Group("/team")

	authGroup.Post("/add", controllers.TeamAdd)
}
