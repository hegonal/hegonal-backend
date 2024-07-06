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

	// check all ongoing incident and update it if reslove
	for _, incident := range incidents {
		if !utils.StringContains(incident.ConfirmLocation, ServerID) {
			continue
		}

		newIncidentTimeline := models.IncidentTimeline{
			IncidentTimelineID: utils.GenerateId(),
			IncidentID:         incident.IncidentID,
			StatusType:         models.IncidentTimelineTypeRecover,
			Message:            "Monitor recover.",
			ServerID:           &ServerID,
			CreatedAt:          utils.TimeNow(),
			UpdatedAt:          utils.TimeNow(),
		}

		if err := db.CreateNewIncidentTimeline(newIncidentTimeline); err != nil {
			log.Error(err)
			continue
		}

		incident.RecoverLocation = append(incident.RecoverLocation, ServerID)
		if utils.UnorderedEqual(incident.RecoverLocation, incident.ConfirmLocation) {
			timeNow := utils.TimeNow()
			incident.IncidentEnd = &timeNow
			incident.IncidentStatus = models.IncidentStatusReslove
			incident.UpdatedAt = timeNow
			if err := db.UpdateIncident(&incident); err != nil {
				log.Error(err)
				continue
			}

			newIncidentTimeline.IncidentTimelineID = utils.GenerateId()
			newIncidentTimeline.Message = "Resolve this incident"
			newIncidentTimeline.StatusType = models.IncidentTimelineTypeResolve
			newIncidentTimeline.ServerID = nil
			if err := db.CreateNewIncidentTimeline(newIncidentTimeline); err != nil {
				log.Error(err)
				continue
			}
		} else {
			incident.UpdatedAt = utils.TimeNow()
			if err := db.UpdateIncident(&incident); err != nil {
				log.Error(err)
				continue
			}
		}
	}
}