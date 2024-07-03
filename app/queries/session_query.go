package queries

import (
	"github.com/hegonal/hegonal-backend/app/models"
	"github.com/hegonal/hegonal-backend/pkg/utils"
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

func (q *SessionQueries) UpdateSession(userID string, oldSession string, newSession string) error {
	query := `UPDATE sessions SET session = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2 AND session = $3`

	result, err := q.Exec(query, newSession, userID, oldSession)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return &NoRowsAffectedError{Message: "no session found to update"}
	}

	return nil
}

func (q *SessionQueries) RotateSession(userID, oldSession string) (string, error) {
	newSession, err := utils.GenerateSessionString()
	if err != nil {
		return "", err
	}

	if err = q.UpdateSession(userID, oldSession, newSession); err != nil {
		return "", err
	}

	return newSession, nil
}
