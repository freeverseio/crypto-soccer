package useractions

import (
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/storage"
)

type UserActionsHistory struct {
	UserActions
	BlockNumber uint64
}

func NewUserActionsHistory(blockNumber uint64, userActions UserActions) *UserActionsHistory {
	history := UserActionsHistory{}
	history.UserActions = userActions
	history.BlockNumber = blockNumber
	return &history
}

func (b UserActionsHistory) ToStorage(tx *sql.Tx) error {
	for _, tactic := range b.UserActions.Tactics {
		th := storage.NewTacticHistory(tactic)
		if err := th.Insert(tx); err != nil {
			return err
		}
	}
	return nil
}
