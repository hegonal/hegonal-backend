package models

import "time"

type TeamAdd struct {
	Name        string `json:"name" validate:"required,max=64"`
	Description string `json:"description" validate:"max=128"`
}

type Team struct {
	ID          string    `db:"id" json:"id" validate:"required,max=20"`
	Name        string    `db:"name" json:"name" validate:"required,max=64"`
	Description string    `db:"description" json:"description" validate:"max=128"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}