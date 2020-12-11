package postgres

import (
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
)

type UniverseService struct {
	DB *sql.DB
}

func (b UniverseService) MarkForDeletion(id string) error {
	return errors.New("not implemented")
}
