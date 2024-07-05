package queries

import (
	"github.com/hegonal/hegonal-backend/app/models"
	"github.com/jmoiron/sqlx"
)

type IncidentQueries struct {
	*sqlx.DB
}

// CreateNewIncident creates a new incident in the database
func (q *IncidentQueries) CreateNewIncident(incident *models.Incident) error {
	query := `INSERT INTO incidents (
		incident_id, team_id, http_monitor_id, confirm_location, recover_location,
		http_status_code, incident_type, incident_status, incident_message, notifications,
		incident_end, incident_start, updated_at
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
	)`

	_, err := q.Exec(
		query,
		incident.IncidentID, incident.TeamID, incident.HttpMonitorID, incident.ConfirmLocation, incident.RecoverLocation,
		incident.HttpStatusCode, incident.IncidentType, incident.IncidentStatus, incident.IncidentMessage, incident.Notifications,
		incident.IncidentEnd, incident.IncidentStart, incident.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

// GetRecentIncident retrieves the most recent incident for a given http_monitor_id and incident_type
func (q *IncidentQueries) GetRecentIncident(httpMonitorID string, incidentType int) (*models.Incident, error) {
	var incident models.Incident
	query := `SELECT
		incident_id, team_id, http_monitor_id, confirm_location, recover_location,
		http_status_code, incident_type, incident_status, incident_message, notifications,
		incident_end, incident_start, updated_at
		FROM incidents
		WHERE http_monitor_id = $1 AND incident_type = $2
		ORDER BY incident_start DESC
		LIMIT 1`

	err := q.Get(&incident, query, httpMonitorID, incidentType)
	if err != nil {
		return nil, err
	}

	return &incident, nil
}

// UpdateIncident updates an existing incident in the database
func (q *IncidentQueries) UpdateIncident(incident *models.Incident) error {
	query := `UPDATE incidents SET
		team_id = $1, http_monitor_id = $2, confirm_location = $3, recover_location = $4,
		http_status_code = $5, incident_type = $6, incident_status = $7, incident_message = $8,
		notifications = $9, incident_end = $10, incident_start = $11, updated_at = $12
		WHERE incident_id = $13`

	_, err := q.Exec(
		query,
		incident.TeamID, incident.HttpMonitorID, incident.ConfirmLocation, incident.RecoverLocation,
		incident.HttpStatusCode, incident.IncidentType, incident.IncidentStatus, incident.IncidentMessage,
		incident.Notifications, incident.IncidentEnd, incident.IncidentStart, incident.UpdatedAt,
		incident.IncidentID,
	)
	if err != nil {
		return err
	}

	return nil
}
