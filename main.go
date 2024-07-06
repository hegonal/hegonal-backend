package main

import (
	"os"

	monitorengine "github.com/hegonal/hegonal-backend/monitor_engine"
	"github.com/hegonal/hegonal-backend/pkg/configs"
	"github.com/hegonal/hegonal-backend/pkg/middleware"
	"github.com/hegonal/hegonal-backend/pkg/routes"
	"github.com/hegonal/hegonal-backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"

	_ "github.com/joho/godotenv/autoload" // load .env file automatically
)

func main() {
	stageStatus := os.Getenv("STAGE_STATUS")

	utils.SnowFlakeInit()

	go monitorengine.RunMonitorEngine()

	if stageStatus == "monitor" {
		log.Info("Monitor mode no running api.")
		select {}
	}

	config := configs.FiberConfig()

	app := fiber.New(config)

	middleware.FiberMiddleware(app)

	route := app.Group("/api")

	routes.AuthRoutes(route)
	routes.TeamRoutes(route)
	routes.NotificationRoutes(route)
	routes.MonitorRoutes(route)

	routes.NotFoundRoute(route)

	if stageStatus == "dev" {
		utils.StartServer(app)
	} else {
		utils.StartServerWithGracefulShutdown(app)
	}
}
