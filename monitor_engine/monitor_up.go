package monitorengine

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/hegonal/hegonal-backend/app/models"
	"github.com/hegonal/hegonal-backend/pkg/utils"
	"github.com/hegonal/hegonal-backend/platform/database"
)

// If incident up check is there have any incident need to be slove
func SloveIncidentHandler(httpMonitor models.HttpMonitor) {
	db, err := database.OpenDBConnection()
	if err != nil {
		log.Error(err)
		return
	}

	incidents, err := db.GetRecentUnsloveIncidents(httpMonitor.HttpMonitorID, models.IncidentStatusOnGoing)
	if err != nil {
		log.Error(err)
		return
	}

	if len(incidents) == 0 {
		return
	}

	for _, incident := range incidents {
		if !utils.StringContains(incident.ConfirmLocation, ServerID) {
			return
		}

		incident.RecoverLocation = append(incident.RecoverLocation, ServerID)
		if utils.UnorderedEqual(incident.RecoverLocation, incident.ConfirmLocation) {
			timeNow := utils.TimeNow()
			incident.IncidentEnd = &timeNow
			incident.IncidentStatus = models.IncidentStatusSlove
			incident.UpdatedAt = timeNow
			err := db.UpdateIncident(&incident)
			if err != nil {
				log.Error(err)
				return
			}
		}
		incident.UpdatedAt = utils.TimeNow()
		err := db.UpdateIncident(&incident)
		if err != nil {
			log.Error(err)
			return
		}
	}

}
