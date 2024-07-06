package queries

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/hegonal/hegonal-backend/app/models"
	"github.com/jmoiron/sqlx"
)

type NotificationQueries struct {
	*sqlx.DB
}

func (q *NotificationQueries) CreateNewNotification(notification *models.Notification) error {
	query := `
	INSERT INTO notifications (
        notification_id,
        user_id,
        team_id,
        notification_type,
        notification_config,
        created_at,
        updated_at
    )
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := q.Exec(
		query,
		notification.NotificationID, notification.UserID, notification.TeamID, notification.NotificationType, notification.NotificationConfig,
		notification.CreatedAt, notification.UpdatedAt,
	)
	return err
}

func (q *NotificationQueries) CreateNewHttpMonitorNotification(notification *models.HttpMonitorNotification) error {
	query := `
	INSERT INTO http_monitor_notifications (http_monitor_id, notification_id)
	VALUES ($1, $2)
	`

	_, err := q.Exec(
		query,
		notification.HttpMonitorID, notification.NotificationID,
	)
	return err
}

func (q *NotificationQueries) GetHttpMonitorNotifications(httpMonitorID string) ([]models.HttpMonitorNotification, error) {
	query := `
	SELECT *
	FROM http_monitor_notifications
	WHERE http_monitor_id = $1
	`

	httpMonitorNotification := []models.HttpMonitorNotification{}

	err := q.Select(&httpMonitorNotification, query, httpMonitorID)
	if err != nil {
		return nil, err
	}

	return httpMonitorNotification, nil
}

func (q *NotificationQueries) GetAllMatchIDNotifications(notificationIDs []string) ([]models.Notification, error) {
	query := `
	SELECT *
	FROM notifications
	WHERE notification_id IN (?)
	`

	query, args, err := sqlx.In(query, notificationIDs)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	query = q.Rebind(query)
	
	var notification = []models.Notification{}
	err = q.Select(&notification, query, args...)

	if err != nil {
		return nil, err
	}

	return notification, nil
}

func (q *NotificationQueries) GetNotificationByID(notificationID string) (*models.Notification, error) {
	query := `
	SELECT *
	FROM notifications
	WHERE notification_id = $1
	`

	var notification = models.Notification{}
	err := q.Get(notification, query, notificationID)

	if err != nil {
		return nil, err
	}

	return &notification, nil
}

func (q *NotificationQueries) GetNotificationsByTeamID(teamID int64) ([]models.Notification, error) {
	query := `
	SELECT *
	FROM notifications
	WHERE team_id = $1
	`

	notifications := []models.Notification{}

	err := q.Select(&notifications, query, teamID)
	if err != nil {
		return nil, err
	}

	return notifications, nil
}

func (q *NotificationQueries) GetNotificationsByUserID(userID int64) ([]models.Notification, error) {
	query := `
	SELECT *
	FROM notifications
	WHERE user_id = $1
	`

	notifications := []models.Notification{}

	err := q.Select(&notifications, query, userID)
	if err != nil {
		return nil, err
	}

	return notifications, nil
}
