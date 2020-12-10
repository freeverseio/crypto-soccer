package postgres

import (
	"database/sql"
	"errors"

	"github.com/freeverseio/crypto-soccer/go/purchasevoider"
	_ "github.com/lib/pq"
)

type UniverseService struct {
	db *sql.DB
}

func NewUniverseService(db *sql.DB) purchasevoider.UniverseService {
	return &UniverseService{
		db: db,
	}
}

func (b UniverseService) MarkForDeletion(id string) error {
	return errors.New("not implemented")
}
