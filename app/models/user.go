package models

import "time"

type UserRole int

const (
	HegonalUser UserRole = iota
	HegonalAdmin
	HegonalOwner
)

type User struct {
	UserID            string    `db:"user_id" json:"user_id" validate:"required,max=20"`
	Name          string    `db:"name" json:"name" validate:"required,max=32"`
	Password      string    `db:"password" json:"password" validate:"required,len=60"`
	Email         string    `db:"email" json:"email" validate:"required,email,max=255"`
	Avatar        string    `db:"avatar" json:"avatar" validate:"max=255"`
	Role          UserRole  `db:"role" json:"role" validate:"oneof=0 1 2"`
	TwoFactorAuth string    `db:"two_factor_auth" json:"two_factor_auth"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}
