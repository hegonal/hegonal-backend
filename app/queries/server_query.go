package queries

import (
	"github.com/hegonal/hegonal-backend/app/models"
	"github.com/jmoiron/sqlx"
)

type ServerQueries struct {
	*sqlx.DB
}

func (q *ServerQueries) CreateNewServerLocation(serverLocation *models.ServerLocation) error {
	query := `
	INSERT INTO server_locations (server_id, server_display_name, country)
	VALUES ($1, $2, $3)
	`

	_, err := q.Exec(
		query,
		serverLocation.ServerID, serverLocation.ServerDisplayName, serverLocation.Country,
	)
	if err != nil {
		return err
	}

	return nil
}

func (q *ServerQueries) ServerLocationExists(serverID string) (bool, error) {
	var exists bool
	query := `
	SELECT EXISTS (
        SELECT 1
        FROM server_locations
        WHERE server_id = $1
    )	
	`

	err := q.Get(&exists, query, serverID)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (q *ServerQueries) GetServerLocation(serverID string) (*models.ServerLocation, error) {
	var serverLocation models.ServerLocation
	query := `
	SELECT server_id,
		server_display_name,
		country
	FROM server_locations
	WHERE server_id = $1
	`

	err := q.Get(&serverLocation, query, serverID)
	if err != nil {
		return nil, err
	}

	return &serverLocation, nil
}

