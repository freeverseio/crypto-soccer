package postgres

import (
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/storage"
	"github.com/freeverseio/crypto-soccer/go/useractions"
)

type UserActionsHistoryStorageService struct {
	UserActionsStorageService
	blockNumber uint64
}

func NewUserActionsHistoryStorageService(tx *sql.Tx, blockNumber uint64) *UserActionsHistoryStorageService {
	return &UserActionsHistoryStorageService{
		UserActionsStorageService: *NewUserActionsStorageService(tx),
		blockNumber:               blockNumber,
	}
}

func (b UserActionsHistoryStorageService) Insert(actions useractions.UserActions) error {
	if err := b.UserActionsStorageService.Insert(actions); err != nil {
		return err
	}
	for _, tactic := range actions.Tactics {
		th := storage.NewTacticHistory(b.blockNumber, tactic)
		if err := th.Insert(b.tx); err != nil {
			return err
		}
	}
	return nil
}
