package queries

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/hegonal/hegonal-backend/app/models"
	"github.com/jmoiron/sqlx"
)

type PingDataQueries struct {
	*sqlx.DB
}

func (q *PingDataQueries) CreatePingData(data *models.PingData) error {
	query := `
		INSERT INTO ping_data (time, http_monitor_id, ping, server_id)
		VALUES (:time, :http_monitor_id, :ping, :server_id)
	`

	_, err := q.NamedExec(query, data)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
