package useractions

import (
	"database/sql"
	"encoding/json"

	"github.com/freeverseio/crypto-soccer/go/relay/storage"
)

type UserActions struct {
	Tactics  []*storage.Tactic   `json:"tactics"`
	Training []*storage.Training `json:"training"`
}

func NewUserActionsByVerse(tx *sql.Tx, verse int) (*UserActions, error) {
	// var err error
	var userActions UserActions
	// if userActions.Tactics, err = storage.TacticsByVerse(verse); err != nil {
	// 	return nil, err
	// }
	return &userActions, nil
}

func (b *UserActions) Marshal() ([]byte, error) {
	if b.Tactics == nil {
		b.Tactics = make([]*storage.Tactic, 0)
	}
	if b.Training == nil {
		b.Training = make([]*storage.Training, 0)
	}
	return json.Marshal(b)
}

func (b *UserActions) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &b)
}
