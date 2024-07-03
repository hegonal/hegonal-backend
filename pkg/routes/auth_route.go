package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hegonal/hegonal-backend/app/controllers"
)

func AuthRoutes(a fiber.Router) {
	authGroup := a.Group("/auth")

	authGroup.Post("/signup", controllers.UserSignUp)
	authGroup.Post("/login", controllers.UserLogin)
	authGroup.Post("/logout", controllers.UserLogout)
}
