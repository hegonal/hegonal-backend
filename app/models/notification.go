package models

import (
	"encoding/json"
	"time"
)

type NotificationType int

const (
	DiscordWebhook NotificationType = iota
)

type Notification struct {
	NotificationID     string           `db:"notification_id" json:"notification_id" validate:"required,max=20"`
	TeamID             *string          `db:"team_id" json:"team_id"`
	UserID             *string          `db:"user_id" json:"user_id"`
	NotificationType   NotificationType `db:"notification_type" json:"notification_type" validate:"oneof=0"`
	NotificationConfig json.RawMessage  `db:"notification_config" json:"notification_config" validate:"required"`
	CreatedAt          time.Time        `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time        `db:"updated_at" json:"updated_at"`
}

type CreateNotification struct {
	TeamID             *string          `db:"team_id" json:"team_id"`
	UserID             *string          `db:"user_id" json:"user_id"`
	NotificationType   NotificationType `db:"notification_type" json:"notification_type" validate:"oneof=0"`
	NotificationConfig json.RawMessage  `db:"notification_config" json:"notification_config" validate:"required"`
}

type HttpMonitorNotification struct {
	HttpMonitorID  string `db:"http_monitor_id" json:"http_monitor_id" validate:"required,max=20"`
	NotificationID string `db:"notification_id" json:"notification_id" validate:"required,max=20"`
}

type CreateHttpMonitorNotification struct {
	TeamID         string `db:"team_id" json:"team_id" validate:"required,max=20"`
	HttpMonitorID  string `db:"http_monitor_id" json:"http_monitor_id" validate:"required,max=20"`
	NotificationID string `db:"notification_id" json:"notification_id" validate:"required,max=20"`
}