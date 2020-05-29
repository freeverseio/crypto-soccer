package postgres

import (
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/storage"
	"github.com/freeverseio/crypto-soccer/go/useractions"
)

type UserActionsStorageService struct {
	tx *sql.Tx
}

func NewUserActionsStorageService(tx *sql.Tx) *UserActionsStorageService {
	return &UserActionsStorageService{
		tx: tx,
	}
}

func (b UserActionsStorageService) UserActionsByTimezone(timezone int) (*useractions.UserActions, error) {
	var err error
	var ua useractions.UserActions
	if ua.Tactics, err = storage.TacticsByTimezone(b.tx, timezone); err != nil {
		return nil, err
	}
	if ua.Trainings, err = storage.TrainingsByTimezone(b.tx, timezone); err != nil {
		return nil, err
	}
	return &ua, nil
}

func (b UserActionsStorageService) Insert(actions useractions.UserActions) error {
	for _, training := range actions.Trainings {
		if err := training.Insert(b.tx); err != nil {
			return err
		}
	}
	for _, tactic := range actions.Tactics {
		if err := tactic.Insert(b.tx); err != nil {
			return err
		}
	}
	return nil
}

// TODO remove this
func (b UserActionsStorageService) InsertHistory(blockNumber uint64, actions useractions.UserActions) error {
	for _, tactic := range actions.Tactics {
		th := storage.NewTacticHistory(blockNumber, tactic)
		if err := th.Insert(b.tx); err != nil {
			return err
		}
	}
	return nil
}
