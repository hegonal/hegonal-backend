package models

type ServerLocation struct {
	ServerID    string  `db:"server_id" json:"server_id" validate:"required,max=32"`
	ServerDisplayName string  `db:"server_display_name" json:"server_display_name" validate:"required,max=32"`
	Country           *string `db:"country" json:"country"`
}
