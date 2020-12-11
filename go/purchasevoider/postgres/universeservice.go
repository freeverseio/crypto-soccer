package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type UniverseService struct {
	DB *sql.DB
}

func (b UniverseService) MarkForDeletion(id string) error {
	_, err := b.DB.Exec("UPDATE INTO players SET voided='true' WHERE player_id=$1", id)
	if err != nil {
		return err
	}
	return nil
}
