package models

import "time"

type Session struct {
	UserID     string    `db:"user_id" json:"user_id" validate:"required,len=18"`
	Session    string    `db:"session" json:"session" validate:"required,len=128"`
	ExpiryTime time.Time `db:"expiry_time" json:"expircy_time" validate:"required"`
	Ip         string    `db:"ip" json:"ip" validate:"required,ip,max=64"`
	Device     string    `db:"device" json:"device" validate:"required,ip,max=255"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}
