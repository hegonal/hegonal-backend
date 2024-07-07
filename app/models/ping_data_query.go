package models

import "time"

type PingData struct {
	Time          time.Time `db:"time" json:"time"`
	HttpMonitorID string    `db:"http_monitor_id" json:"http_monitor_id"`
	Ping          int       `db:"ping" json:"ping"`
	ServerID     string    `db:"server_id" json:"server_id"`
}
