package queries

import (
	"time"

	"github.com/hegonal/hegonal-backend/app/models"
	"github.com/hegonal/hegonal-backend/pkg/utils"
	"github.com/jmoiron/sqlx"
)

type SessionQueries struct {
	*sqlx.DB
}

func (q *SessionQueries) CreateNewSession(session *models.Session) error {
	query := `INSERT INTO sessions VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := q.Exec(
		query,
		session.UserID, session.Session, session.ExpiryTime, session.Ip, session.Device, session.CreatedAt, session.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (q *UserQueries) GetSession(userID, userSession string) (models.Session, error) {
	var session models.Session
	query := `
	SELECT user_id,
		session,
		expiry_time,
		ip,
		device,
		created_at,
		updated_at
	FROM sessions
	WHERE user_id = $1
		AND session = $2
	`
	err := q.Get(&session, query, userID, userSession)
	if err != nil {
		return session, err
	}
	return session, nil
}

func (q *SessionQueries) DeleteSession(userID string, session string) error {
	query := `
	DELETE FROM sessions
	WHERE user_id = $1
    	AND session = $2
	`

	_, err := q.Exec(query, userID, session)
	if err != nil {
		return err
	}

	return nil
}

func (q *SessionQueries) UpdateSession(userID string, oldSession string, newSession string, expriceTime time.Time) error {
	query := `
	UPDATE sessions
	SET session = $1,
		expiry_time = $2,
		updated_at = CURRENT_TIMESTAMP
	WHERE user_id = $3
		AND session = $4
	`

	result, err := q.Exec(query, newSession, expriceTime, userID, oldSession)
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

	if err = q.UpdateSession(userID, oldSession, newSession, utils.TimeNow().Add(24*time.Hour)); err != nil {
		return "", err
	}

	return newSession, nil
}
