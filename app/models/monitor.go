package models

import (
	"time"

	"github.com/lib/pq"
)

type HttpMethod int

const (
	HttpGet HttpMethod = iota
	HttpPost
	HttpPut
	HttpPatch
	HttpDelete
	HttpHead
	HttpOptions
)

type BodyEncoding int

const (
	BodyText BodyEncoding = iota
	BodyJavaScript
	BodyJSON
	BodyHTML
	BodyXML
)

type HttpMonitorStatus int

const (
	HttpMonitorStatusStop HttpMonitorStatus = iota
	HttpMonitorStatusUp
	HttpMonitorStatusDown
	HttpMonitorStatusUnknow
)

type HttpMonitor struct {
	HttpMonitorID      string            `db:"http_monitor_id" json:"http_monitor_id" validate:"required,max=20"`
	TeamID             string            `db:"team_id" json:"team_id" validate:"required,max=20"`
	Status             HttpMonitorStatus `db:"status" json:"status" validate:"required"`
	URL                string            `db:"url" json:"url" validate:"required,http_url,max=255"`
	Interval           int               `db:"interval" json:"interval" validate:"required,lte=86400"`
	Retries            int               `db:"retries" json:"retries" validate:"required,lte=32767"`
	RetryInterval      int               `db:"retry_interval" json:"retry_interval" validate:"required,lte=86400"`
	RequestTimeout     int               `db:"request_timeout" json:"request_timeout" validate:"required,lte=60"`
	ResendNotification int               `db:"resend_notification" json:"resend_notification" validate:"required,lte=32767"`
	FollowRedirections bool              `db:"follow_redirections" json:"follow_redirections" validate:"required"`
	MaxRedirects       int               `db:"max_redirects" json:"max_redirects" validate:"required,lte=32767"`
	CheckSslError      bool              `db:"check_ssl_error" json:"check_ssl_error" validate:"required"`
	// day
	SslExpiryReminders int `db:"ssl_expiry_reminders" json:"ssl_expiry_reminders" validate:"required,lte=32767"`
	// day
	DomainExpiryReminders int            `db:"domain_expiry_reminders" json:"domain_expiry_reminders" validate:"required,lte=32767"`
	HttpStatusCodes       pq.StringArray `db:"http_status_codes" json:"http_status_codes" validate:"required,dive,max=3"`
	HttpMethod            HttpMethod     `db:"http_method" json:"http_method" validate:"oneof=0 1 2 3 4 5 6"`
	BodyEncoding          *BodyEncoding  `db:"body_encoding" json:"body_encoding"`
	RequestBody           *string        `db:"request_body" json:"request_body"`
	RequestHeaders        *string        `db:"request_headers" json:"request_headers"`
	Group                 *string        `db:"group" json:"group"`
	Proxy                 *string        `db:"proxy" json:"proxy"`
	SendToOnCall          bool           `db:"send_to_oncall" json:"send_to_oncall"`
	CreatedAt             time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt             time.Time      `db:"updated_at" json:"updated_at"`
}

type CreateNewHttpMonitor struct {
	TeamID             string `db:"team_id" json:"team_id" validate:"required,max=20"`
	URL                string `db:"url" json:"url" validate:"required,http_url,max=255"`
	Interval           int    `db:"interval" json:"interval" validate:"required,lte=86400"`
	Retries            int    `db:"retries" json:"retries" validate:"required,lte=32767"`
	RetryInterval      int    `db:"retry_interval" json:"retry_interval" validate:"required,lte=86400"`
	RequestTimeout     int    `db:"request_timeout" json:"request_timeout" validate:"required,lte=60"`
	ResendNotification int    `db:"resend_notification" json:"resend_notification" validate:"required,lte=32767"`
	FollowRedirections bool   `db:"follow_redirections" json:"follow_redirections" validate:"required"`
	MaxRedirects       int    `db:"max_redirects" json:"max_redirects" validate:"required,lte=32767"`
	CheckSslError      bool   `db:"check_ssl_error" json:"check_ssl_error" validate:"required"`
	// day
	SslExpiryReminders int `db:"ssl_expiry_reminders" json:"ssl_expiry_reminders" validate:"required,lte=32767"`
	// day
	DomainExpiryReminders int            `db:"domain_expiry_reminders" json:"domain_expiry_reminders" validate:"required,lte=32767"`
	HttpStatusCodes       pq.StringArray `db:"http_status_codes" json:"http_status_codes" validate:"required,dive,max=3"`
	HttpMethod            HttpMethod     `db:"http_method" json:"http_method" validate:"oneof=0 1 2 3 4 5 6"`
	BodyEncoding          *BodyEncoding  `db:"body_encoding" json:"body_encoding"`
	RequestBody           *string        `db:"request_body" json:"request_body"`
	RequestHeaders        *string        `db:"request_headers" json:"request_headers"`
	Group                 *string        `db:"group" json:"group"`
	Proxy                 *string        `db:"proxy" json:"proxy"`
	SendToOnCall          bool           `db:"send_to_oncall" json:"send_to_oncall"`
}