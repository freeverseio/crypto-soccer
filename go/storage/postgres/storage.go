package postgres

import (
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/storage"
)

func NewStorageService(tx *sql.Tx) *storage.StorageService {
	return &storage.StorageService{
		TeamService: NewTeamStorageService(tx),
	}
}
