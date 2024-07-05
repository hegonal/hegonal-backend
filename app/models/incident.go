package models

import (
	"time"

	"github.com/lib/pq"
)

type IncidentType int

const (
	IncidentTypeStatusCodeNotMatch IncidentType = iota
	IncidentTypeTimeout
	IncidentTypeOtherError
	IncidentTypeSSLError
	IncidentTypeSSLExpiry
	IncidentTypeDomainExpiry
	IncidentTypeKeyWordNotFound
	IncidentTypeFine
)

type Incident struct {
	IncidentID      string         `db:"incident_id" json:"incident_id" validate:"required,max=20"`
	TeamID          string         `db:"team_id" json:"team_id" validate:"required,max=20"`
	HttpMonitorID   *string        `db:"http_monitor_id" json:"http_monitor_id"`
	HttpStatusCode  *string        `db:"http_status_code" json:"http_status_code"`
	ConfirmLocation pq.StringArray `db:"confirm_location" json:"confirm_location"`
	RecoverLocation pq.StringArray `db:"recover_location" json:"recover_location"`
	IncidentType    int            `db:"incident_type" json:"incident_type" validate:"required"`
	IncidentStatus  int16          `db:"incident_status" json:"incident_status" validate:"required"`
	IncidentMessage string         `db:"incident_message" json:"incident_message" validate:"required"`
	Notifications   bool           `db:"notifications" json:"notifications" validate:"required"`
	IncidentEnd     time.Time      `db:"incident_end" json:"incident_end" validate:"required"`
	IncidentStart   time.Time      `db:"incident_start" json:"incident_start" validate:"required"`
	UpdatedAt       time.Time      `db:"updated_at" json:"updated_at" validate:"required"`
}
