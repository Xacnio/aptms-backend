package database

import "github.com/Xacnio/aptms-backend/app/queries"

type Queries struct {
	*queries.UserQueries
	*queries.AptQueries
	*queries.InitQueries
	*queries.CDistQueries
}

func OpenDBConnection() (*Queries, error) {
	db, err := PostgreSQLConnection()
	if err != nil {
		return nil, err
	}

	return &Queries{
		UserQueries:  &queries.UserQueries{DB: db},
		AptQueries:   &queries.AptQueries{DB: db},
		InitQueries:  &queries.InitQueries{DB: db},
		CDistQueries: &queries.CDistQueries{DB: db},
	}, nil
}
