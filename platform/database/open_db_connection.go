package database

import (
	"github.com/hegonal/hegonal-backend/app/queries"
	"github.com/jmoiron/sqlx"
)

type Queries struct {
	*queries.UserQueries
	*queries.SessionQueries
	*queries.TeamQueries
	*queries.MonitorQueries
}

func OpenDBConnection() (*Queries, error) {
	var (
		db  *sqlx.DB
		err error
	)

	db, err = PostgreSQLConnection()

	if err != nil {
		return nil, err
	}

	return &Queries{
		UserQueries:    &queries.UserQueries{DB: db},
		SessionQueries: &queries.SessionQueries{DB: db},
		TeamQueries:    &queries.TeamQueries{DB: db},
		MonitorQueries: &queries.MonitorQueries{DB: db},
	}, nil
}
