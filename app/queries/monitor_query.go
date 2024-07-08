package queries

import (
	"github.com/hegonal/hegonal-backend/app/models"
	"github.com/jmoiron/sqlx"
)

type MonitorQueries struct {
	*sqlx.DB
}

func (q *MonitorQueries) CreateNewHttpMonitor(httpMonitor *models.HttpMonitor) error {
	query := `
	INSERT INTO http_monitors (
			http_monitor_id,
			team_id,
			status,
			url,
			interval,
			retries,
			retry_interval,
			request_timeout,
			resend_notification,
			follow_redirections,
			max_redirects,
			check_ssl_error,
			ssl_expiry_reminders,
			domain_expiry_reminders,
			http_status_codes,
			http_method,
			body_encoding,
			request_body,
			request_headers,
			"group",
			proxy,
			send_to_oncall,
			created_at,
			updated_at
		)
	VALUES (
			:http_monitor_id,
			:team_id,
			:status,
			:url,
			:interval,
			:retries,
			:retry_interval,
			:request_timeout,
			:resend_notification,
			:follow_redirections,
			:max_redirects,
			:check_ssl_error,
			:ssl_expiry_reminders,
			:domain_expiry_reminders,
			:http_status_codes,
			:http_method,
			:body_encoding,
			:request_body,
			:request_headers,
			:group,
			:proxy,
			:send_to_oncall,
			:created_at,
			:updated_at
		)
	`

	_, err := q.DB.NamedExec(query, httpMonitor)
	if err != nil {
		return err
	}

	return nil
}

func (q *MonitorQueries) GetHttpMonitorByID(httpMonitorID string) (*models.HttpMonitor, error) {
	query := `
		SELECT * FROM http_monitors WHERE http_monitor_id = $1
	`

	var httpMonitor models.HttpMonitor
	err := q.DB.Get(&httpMonitor, query, httpMonitorID)
	if err != nil {
		return nil, err
	}

	return &httpMonitor, nil
}

func (q *MonitorQueries) GetAllHttpMonitors() ([]models.HttpMonitor, error) {
	query := `
		SELECT * FROM http_monitors
	`

	var httpMonitors []models.HttpMonitor
	err := q.DB.Select(&httpMonitors, query)
	if err != nil {
		return nil, err
	}

return httpMonitors, nil
}

func (q *MonitorQueries) GetAllTeamHttpMonitors(teamID string) ([]models.HttpMonitor, error) {
	query := `
	SELECT *
	FROM http_monitors
	WHERE team_id = $1;
	`

	var httpMonitors []models.HttpMonitor
	err := q.Select(&httpMonitors, query, teamID)
	if err != nil {
		return nil, err
	}

	return httpMonitors, nil
}

func (q *MonitorQueries) DeleteHttpMonitor(httpMonitorID string) error {
	query := `
	DELETE FROM http_monitors
	WHERE http_monitor_id = $1;
	`

	_, err := q.Exec(query, httpMonitorID)
	if err != nil {
		return err
	}

	return nil
}