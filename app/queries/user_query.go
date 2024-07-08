package queries

import (
	"github.com/hegonal/hegonal-backend/app/models"
	"github.com/jmoiron/sqlx"
)

type UserQueries struct {
	*sqlx.DB
}

func (q *UserQueries) IsEmailUsed(email string) (bool, error) {
	var count int
	query := `
	SELECT COUNT(*)
	FROM users
	WHERE email = $1
	`
	err := q.Get(&count, query, email)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (q *UserQueries) IsOwnerAccountCreated() (bool, error) {
	var exists bool
	query := `
	SELECT EXISTS(
        SELECT 1
        FROM users
        WHERE role = $1
    )
	`
	err := q.Get(&exists, query, 0)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (q *UserQueries) CreateNewUser(u *models.User) error {
	query := `
	INSERT INTO users
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := q.Exec(
		query,
		u.UserID, u.Name, u.Password, u.Email, u.Avatar, u.Role, u.TwoFactorAuth, u.CreatedAt, u.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (q *UserQueries) GetUserByEmail(email string) (models.User, error) {
	var u models.User
	query := `
	SELECT *
	FROM users
	WHERE email = $1
	`

	err := q.Get(&u, query, email)
	if err != nil {
		return u, err
	}

	return u, nil
}
