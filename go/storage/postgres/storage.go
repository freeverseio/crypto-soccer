package postgres

import (
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/storage"
)

type StorageService struct {
	TeamStorageService  TeamStorageService
	MatchStorageService MatchStorageService
}

func NewStorageService(tx *sql.Tx) *StorageService {
	return &StorageService{
		*NewTeamStorageService(tx),
		*NewMatchStorageService(tx),
	}
}

func (b StorageService) TeamService() storage.TeamStorageService {
	return b.TeamStorageService
}
func (b StorageService) MatchService() storage.MatchStorageService {
	return b.MatchStorageService
}
