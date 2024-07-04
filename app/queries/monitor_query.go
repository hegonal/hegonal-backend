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
		INSERT INTO http_monitors (http_monitor_id, team_id, active, url, interval, retries, retry_interval,
			request_timeout, resend_notification, follow_redirections, max_redirects, check_ssl_error,
			ssl_expiry_reminders, domain_expiry_reminders, http_status_codes, http_method, body_encoding,
			request_body, request_headers, "group", notification, proxy, created_at, updated_at)
		VALUES (:http_monitor_id, :team_id, :active, :url, :interval, :retries, :retry_interval,
			:request_timeout, :resend_notification, :follow_redirections, :max_redirects, :check_ssl_error,
			:ssl_expiry_reminders, :domain_expiry_reminders, :http_status_codes, :http_method, :body_encoding,
			:request_body, :request_headers, :group, :notification, :proxy, :created_at, :updated_at)
	`

	_, err := q.DB.NamedExec(query, httpMonitor)
	if err != nil {
		return err
	}

	return nil
}
