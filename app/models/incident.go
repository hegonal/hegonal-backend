package models

import (
	"time"

	"github.com/lib/pq"
)

type IncidentType int

const (
	IncidentTypeFine IncidentType = iota
	IncidentTypeSSLExpiry
	IncidentTypeDomainExpiry
	IncidentTypeKeyWordNotFound
	IncidentTypeStatusCodeNotMatch
	IncidentTypeTimeout
	IncidentTypeOtherError
)

type IncidentStatus int

const (
	IncidentStatusReslove IncidentStatus = iota
	IncidentStatusOnGoing
)

type Incident struct {
	IncidentID      string         `db:"incident_id" json:"incident_id" validate:"required,max=20"`
	TeamID          string         `db:"team_id" json:"team_id" validate:"required,max=20"`
	HttpMonitorID   *string        `db:"http_monitor_id" json:"http_monitor_id"`
	HttpStatusCode  *string        `db:"http_status_code" json:"http_status_code"`
	ExpiryDate      *time.Time     `db:"expiry_date" json:"expiry_date"`
	ConfirmLocation pq.StringArray `db:"confirm_location" json:"confirm_location"`
	RecoverLocation pq.StringArray `db:"recover_location" json:"recover_location"`
	IncidentType    IncidentType   `db:"incident_type" json:"incident_type" validate:"required"`
	IncidentStatus  IncidentStatus `db:"incident_status" json:"incident_status" validate:"required"`
	IncidentMessage string         `db:"incident_message" json:"incident_message" validate:"required"`
	Notifications   bool           `db:"notifications" json:"notifications" validate:"required"`
	IncidentEnd     *time.Time     `db:"incident_end" json:"incident_end"`
	IncidentStart   time.Time      `db:"incident_start" json:"incident_start" validate:"required"`
	UpdatedAt       time.Time      `db:"updated_at" json:"updated_at" validate:"required"`
}

type IncidentTimelineType int

const (
	IncidentTimelineTypeConfirm IncidentTimelineType = iota
	IncidentTimelineTypeRecover
	IncidentTimelineTypeSendNotification
	IncidentTImelineTypeInvestigating
	IncidentTImelineTypeIdentified
	IncidentTImelineTypeUpdate
	IncidentTimelineTypeResolve
)

type IncidentTimeline struct {
	IncidentTimelineID string    `db:"incident_timeline_id" json:"incident_timeline_id" validate:"required,max=20"`
	IncidentID         string    `db:"incident_id" json:"incident_id" validate:"required,max=20"`
	StatusType         IncidentTimelineType       `db:"status_type" json:"status_type"`
	Message            string    `db:"message" json:"message"`
	CreatedBy          *string   `db:"created_by" json:"created_by"`
	ServerID           *string   `db:"server_id" json:"server_id"`
	CreatedAt          time.Time `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time `db:"updated_at" json:"update_at"`
}