package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hegonal/hegonal-backend/app/controllers"
	"github.com/hegonal/hegonal-backend/pkg/middleware"
)

func MonitorRoutes(a fiber.Router) {
	monitorGroup := a.Group("/monitor")

	monitorGroup.Post("/http/add", middleware.SessionValidationMiddleware, controllers.CreateNewHttpMonitor)
	
	monitorGroup.Get("/http/get/:teamID", middleware.SessionValidationMiddleware, controllers.GetAllHttpMonitors)
	monitorGroup.Delete("/http/delete/:teamID/:monitorID", middleware.SessionValidationMiddleware, controllers.DeleteHttpMonitor)
}
