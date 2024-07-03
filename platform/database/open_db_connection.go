package database

import (
	"github.com/hegonal/hegonal-backend/app/queries"
	"github.com/jmoiron/sqlx"
)

type Queries struct {
	*queries.SessionQueries
	*queries.UserQueries
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
		SessionQueries: &queries.SessionQueries{DB: db},
		UserQueries:    &queries.UserQueries{DB: db},
	}, nil
}
