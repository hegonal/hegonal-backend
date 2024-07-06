package monitorengine

import (
	"os"

	"github.com/gofiber/fiber/v2/log"
	"github.com/hegonal/hegonal-backend/app/models"
	"github.com/hegonal/hegonal-backend/platform/database"
)

var ServerID string

func RunMonitorEngine() {
	db, err := database.OpenDBConnection()
	if err != nil {
		log.Error(err)
		panic("error open db connection when start monitor engine")
	}
	
	checkAndCraeteServerLocation(db)

	// Start htt monitor
	httpMonitors, err := db.GetAllHttpMonitors()
	if err != nil {
		log.Error(err)
		panic("error get http monitor data")
	}

	for _, httpMonitor := range httpMonitors {
		if httpMonitor.Status != models.HttpMonitorStatusStop {
			go startHttpMonitor(httpMonitor)
		}
	}

	// other monitor start
}

func checkAndCraeteServerLocation(db *database.Queries) {
	serverID := os.Getenv("SERVER_ID")
	ServerID = serverID
	isServerLocationExist, err := db.ServerLocationExists(serverID)
	if err != nil {
		log.Error(err)
		panic(err.Error())
	}
	if !isServerLocationExist {
		log.Info("Server location not exist,create one.")
		serverLocaton := &models.ServerLocation{}
		serverLocaton.ServerID = serverID
		serverLocaton.ServerDisplayName = serverID
		if err := db.CreateNewServerLocation(serverLocaton); err != nil {
			log.Error(err)
			panic(err.Error())
		}
	}
}
