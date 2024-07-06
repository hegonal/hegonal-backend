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
	NotificationType   NotificationType `db:"notification_type" json:"notification_type" validate:"oneof=0"`
	NotificationConfig json.RawMessage  `db:"notification_config" json:"notification_config" validate:"required"`
	CreatedAt          time.Time        `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time        `db:"updated_at" json:"updated_at"`
}

type CreateNotification struct {
	NotificationType   NotificationType `db:"notification_type" json:"notification_type" validate:"oneof=0"`
	NotificationConfig json.RawMessage  `db:"notification_config" json:"notification_config" validate:"required"`
}
