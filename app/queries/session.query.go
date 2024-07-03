package queries

import (
	"github.com/hegonal/hegonal-backend/app/models"
	"github.com/jmoiron/sqlx"
)

type SessionQueries struct {
	*sqlx.DB
}

func (q *SessionQueries) CreateNewSession(s *models.Session) error {
	query := `INSERT INTO sessions VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := q.Exec(
		query,
		s.ID, s.Session, s.Ip, s.Device, s.CreatedAt, s.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (q *SessionQueries) DeleteSession(userID string, session string) error {
	query := `DELETE FROM sessions WHERE id = $1 AND session = $2`

	_, err := q.Exec(query, userID, session)
	if err != nil {
		return err
	}

	return nil
}
