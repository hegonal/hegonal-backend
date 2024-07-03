package models

type SignUp struct {
	Name          string    `db:"name" json:"name" validate:"required,max=255"`
	Password      string    `db:"password" json:"password" validate:"required,max=72"`
	Email         string    `db:"email" json:"email" validate:"required,email,max=255"`
}

type Login struct {
	Password      string    `db:"password" json:"password" validate:"required,max=72"`
	Email         string    `db:"email" json:"email" validate:"required,email,max=255"`
}
