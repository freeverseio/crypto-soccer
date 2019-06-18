package postgres

import "database/sql"

type PostgresStorage struct {
	db *sql.DB
}

func New(url string) (*PostgresStorage, error) {
	var err error
	storage := &PostgresStorage{}
	storage.db, err = sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	return storage, nil
}
