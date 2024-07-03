package models

import "time"

type Session struct {
	ID        string    `db:"id" json:"id" validate:"required,len=18"`
	Session   string    `db:"session" json:"session" validate:"required,len=128"`
	Ip        string    `db:"ip" json:"ip" validate:"required,ip,max=64"`
	Device    string    `db:"Device" json:"Device" validate:"required,ip,max=255"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"update_at" json:"updated_at"`
}
