package queries

import (
	"github.com/hegonal/hegonal-backend/app/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type IncidentQueries struct {
	*sqlx.DB
}

// CreateNewIncident creates a new incident in the database
func (q *IncidentQueries) CreateNewIncident(incident *models.Incident) error {
	query := `
	INSERT INTO incidents (
			incident_id,
			team_id,
			http_monitor_id,
			expiry_date,
			confirm_location,
			recover_location,
			http_status_code,
			incident_type,
			incident_status,
			incident_message,
			notifications,
			incident_end,
			incident_start,
			updated_at
		)
	VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9,
			$10,
			$11,
			$12,
			$13,
			$14
		)
	`

	_, err := q.Exec(
		query,
		incident.IncidentID, incident.TeamID, incident.HttpMonitorID, incident.ExpiryDate, incident.ConfirmLocation, incident.RecoverLocation,
		incident.HttpStatusCode, incident.IncidentType, incident.IncidentStatus, incident.IncidentMessage, incident.Notifications,
		incident.IncidentEnd, incident.IncidentStart, incident.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

// GetRecentIncident retrieves the most recent incident for a given http_monitor_id and incident_type
func (q *IncidentQueries) GetRecentIncidentWithSpecifyTypeAndStatus(httpMonitorID string, incidentType models.IncidentType, incidentStatus models.IncidentStatus) (*models.Incident, error) {
	var incident models.Incident
	query := `
	SELECT *
	FROM incidents
	WHERE http_monitor_id = $1
		AND incident_type = $2
		AND incident_status = $3
	ORDER BY incident_start DESC
	LIMIT 1
	`

	err := q.Get(&incident, query, httpMonitorID, incidentType, incidentStatus)
	if err != nil {
		return nil, err
	}

	return &incident, nil
}

// UpdateIncident updates an existing incident in the database
func (q *IncidentQueries) UpdateIncident(incident *models.Incident) error {
	query := `
	UPDATE incidents
	SET team_id = $1,
		http_monitor_id = $2,
		expiry_date = $3,
		confirm_location = $4,
		recover_location = $5,
		http_status_code = $6,
		incident_type = $7,
		incident_status = $8,
		incident_message = $9,
		notifications = $10,
		incident_end = $11,
		incident_start = $12,
		updated_at = $13
	WHERE incident_id = $14
	`

	_, err := q.Exec(
		query,
		incident.TeamID, incident.HttpMonitorID, incident.ExpiryDate, incident.ConfirmLocation, incident.RecoverLocation,
		incident.HttpStatusCode, incident.IncidentType, incident.IncidentStatus, incident.IncidentMessage,
		incident.Notifications, incident.IncidentEnd, incident.IncidentStart, incident.UpdatedAt,
		incident.IncidentID,
	)
	if err != nil {
		return err
	}

	return nil
}

// UpdateConfirmLocation updates the ConfirmLocation of an incident in the database
func (q *IncidentQueries) UpdateIncidentConfirmLocation(incidentID string, confirmLocation pq.StringArray) error {
	query := `
	UPDATE incidents
	SET confirm_location = $1,
		updated_at = NOW()
	WHERE incident_id = $2
	`

	_, err := q.Exec(query, confirmLocation, incidentID)
	if err != nil {
		return err
	}

	return nil
}

// GetRecentIncidents retrieves the most recent incidents for a given http_monitor_id and multiple incident_types
func (q *IncidentQueries) GetRecentUnsloveIncidents(httpMonitorID string, incidentStatus models.IncidentStatus) ([]models.Incident, error) {
	var incidents []models.Incident
	query := `
	SELECT *
	FROM incidents
	WHERE http_monitor_id = $1
		AND incident_status = $2
	ORDER BY incident_start DESC
	`

	err := q.Select(&incidents, query, httpMonitorID, incidentStatus)
	if err != nil {
		return nil, err
	}

	return incidents, nil
}

// CreateNewIncidentTimeline creates a new incident timeline entry in the database
func (q *IncidentQueries) CreateNewIncidentTimeline(incidentTimeline models.IncidentTimeline) error {
	query := `
	INSERT INTO incident_timelines (
			incident_timeline_id,
			incident_id,
			status_type,
			message,
			created_by,
			server_id,
			created_at,
			updated_at
		)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := q.Exec(
		query,
		incidentTimeline.IncidentTimelineID, incidentTimeline.IncidentID, incidentTimeline.StatusType, incidentTimeline.Message,
		incidentTimeline.CreatedBy, incidentTimeline.ServerID, incidentTimeline.CreatedAt, incidentTimeline.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}