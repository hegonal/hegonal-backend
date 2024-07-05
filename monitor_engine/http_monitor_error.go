package monitorengine

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/hegonal/hegonal-backend/app/models"
)

// If have error check what this error is and create new incident & send notification
func incidentHandler(httpMonitor *models.HttpMonitor, errorType models.IncidentType, expiry int, httpCode int, err error) {
	log.Info(httpMonitor.URL)
	log.Info(errorType)
	log.Info(expiry)
	log.Info(httpCode)
	log.Info(err.Error())
}