package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hegonal/hegonal-backend/app/controllers"
	"github.com/hegonal/hegonal-backend/pkg/middleware"
)

func MonitorRoutes(a fiber.Router) {
	authGroup := a.Group("/monitor")

	authGroup.Post("/http/add", middleware.SessionValidationMiddleware, controllers.CreateNewHttpMonitor)
}
