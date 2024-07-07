package models

import "time"

type CreateNewTeam struct {
	Name        string `json:"name" validate:"required,max=64"`
	Description string `json:"description" validate:"max=128"`
}

type Team struct {
	TeamID      string    `db:"team_id" json:"team_id" validate:"required,max=20"`
	Name        string    `db:"name" json:"name" validate:"required,max=64"`
	Description string    `db:"description" json:"description" validate:"max=128"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

type TeamMemberRole int

const (
	TeamNoneRole TeamMemberRole = iota
	TeamViewers
	TeamUser
	TeamAdmin
	TeamOwner
)

type TeamMember struct {
	MemberID  string         `db:"member_id" json:"member_id" validate:"required,max=20"`
	TeamID    string         `db:"team_id" json:"team_id" validate:"required,max=20"`
	Role      TeamMemberRole `db:"role" json:"role" validate:"oneof=0 1 2 3"`
	CreatedAt time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt time.Time      `db:"updated_at" json:"updated_at"`
}

type TeamInvite struct {
	InvitedID string `db:"invited_id" json:"invited_id" validate:"required,max=20"`
	TeamID    string `db:"team_id" json:"team_id" validate:"required,max=20"`
	UserID    string `db:"user_id" json:"user_id" validate:"required,max=20"`
	ExpiryDate time.Time `db:"expiry_date" json:"expiry_date"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type CreateNewTeamInvite struct {
	TeamID    string `db:"team_id" json:"team_id" validate:"required,max=20"`
	Email    string `db:"email" json:"email" validate:"required,max=255"`
}
