package main

import (
	"os"

	"github.com/hegonal/hegonal-backend/pkg/configs"
	"github.com/hegonal/hegonal-backend/pkg/middleware"
	"github.com/hegonal/hegonal-backend/pkg/routes"
	"github.com/hegonal/hegonal-backend/pkg/utils"

	"github.com/gofiber/fiber/v2"

	_ "github.com/joho/godotenv/autoload" // load .env file automatically
)

func main() {
	config := configs.FiberConfig()

	app := fiber.New(config)

	middleware.FiberMiddleware(app)

	utils.SnowFlakeInit()

	route := app.Group("/api")

	routes.AuthRoutes(route)
	routes.TeamRoutes(route)
	routes.NotFoundRoute(route)

	if os.Getenv("STAGE_STATUS") == "dev" {
		utils.StartServer(app)
	} else {
		utils.StartServerWithGracefulShutdown(app)
	}
}
