package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hegonal/hegonal-backend/app/controllers"
	"github.com/hegonal/hegonal-backend/pkg/middleware"
)

func NotificationRoutes(a fiber.Router) {
	authGroup := a.Group("/notification")

	authGroup.Post("/add", middleware.SessionValidationMiddleware, controllers.CreateNewNotification)
	authGroup.Post("/monitor/http/add", middleware.SessionValidationMiddleware, controllers.CreateNewHttpMonitorNotification)
}
