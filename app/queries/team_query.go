package queries

import "github.com/hegonal/hegonal-backend/app/models"

func (q *SessionQueries) CreateNewTeam(s *models.Team) error {
	query := `INSERT INTO teams VALUES ($1, $2, $3, $4, $5)`

	_, err := q.Exec(
		query,
		s.ID, s.Name, s.Description, s.CreatedAt, s.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (q *SessionQueries) MonitorAdd(s *models.Team) error {
	query := `INSERT INTO teams VALUES ($1, $2, $3, $4, $5)`

	_, err := q.Exec(
		query,
		s.ID, s.Name, s.Description, s.CreatedAt, s.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}