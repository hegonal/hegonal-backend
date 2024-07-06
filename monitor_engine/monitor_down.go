package monitorengine

import (
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/hegonal/hegonal-backend/app/models"
	"github.com/hegonal/hegonal-backend/pkg/utils"
	"github.com/hegonal/hegonal-backend/platform/database"
)

// If have error check what this error is and create new incident & send notification
func incidentHandler(httpMonitor models.HttpMonitor, errorType models.IncidentType, expiry int, httpCode int, incidentErr error) {
	db, err := database.OpenDBConnection()
	if err != nil {
		log.Error(err)
		return
	}

	var newIncident = &models.Incident{}

	newIncident.IncidentID = utils.GenerateId()
	newIncident.TeamID = httpMonitor.TeamID
	newIncident.HttpMonitorID = &httpMonitor.HttpMonitorID
	newIncident.IncidentStatus = models.IncidentStatusOnGoing
	if httpCode == 0 {
		newIncident.HttpStatusCode = nil
	} else {
		httpStatusCodePointer := new(string)
		*httpStatusCodePointer = string(httpCode)
		newIncident.HttpStatusCode = httpStatusCodePointer
	}
	if expiry == 0 {
		newIncident.ExpiryDate = nil
	} else {
		exiryDate := new(time.Time)
		*exiryDate = utils.TimeNow().Add(time.Duration(expiry) * time.Hour * 24)
		newIncident.ExpiryDate = exiryDate
	}
	newIncident.ConfirmLocation = []string{ServerID}
	newIncident.RecoverLocation = []string{}
	newIncident.IncidentType = errorType
	if incidentErr != nil {
		newIncident.IncidentMessage = incidentErr.Error()
	}
	newIncident.Notifications = true
	newIncident.IncidentEnd = nil
	newIncident.IncidentStart = utils.TimeNow()
	newIncident.UpdatedAt = utils.TimeNow()
	
	checkOrUpdateIncident(newIncident, &httpMonitor, db, errorType)
}

func checkOrUpdateIncident(newIncident *models.Incident, httpMonitor *models.HttpMonitor, db *database.Queries, errorType models.IncidentType) {
	newIncident.IncidentStatus = models.IncidentStatusOnGoing

	incident, err := db.GetRecentIncident(httpMonitor.HttpMonitorID, errorType)
	if err == sql.ErrNoRows {
		err := db.CreateNewIncident(newIncident)
		if err != nil {
			log.Error(err)
			return
		}
		return
	} else if err != nil {
		log.Error(err)
		return
	}

	if incident.IncidentStatus != models.IncidentStatusOnGoing {
		err := db.CreateNewIncident(newIncident)
		if err != nil {
			log.Error(err)
			return
		}
		return
	}

	if utils.StringContains(incident.ConfirmLocation, ServerID) {
		return
	} else {
		incident.ConfirmLocation = append(incident.ConfirmLocation, ServerID)
		err := db.UpdateIncidentConfirmLocation(incident.IncidentID, incident.ConfirmLocation)
		if err != nil {
			log.Error(err)
			return
		}
	}
}
