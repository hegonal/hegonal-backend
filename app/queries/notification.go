package queries

import (
	"github.com/hegonal/hegonal-backend/app/models"
	"github.com/jmoiron/sqlx"
)

type NotificationQueries struct {
	*sqlx.DB
}

func (q *NotificationQueries) CreateNewNotification(notification *models.Notification) error {
	query := `INSERT INTO notifications (
		notification_id, notification_type, notification_config, created_at, updated_at
	) VALUES (
		$1, $2, $3, $4, $5
	)`

	_, err := q.Exec(
		query,
		notification.NotificationID, notification.NotificationType, notification.NotificationConfig,
		notification.CreatedAt, notification.UpdatedAt,
	)
	return err
}