package storage

import (
	"encoding/json"
)

type UserActions struct {
	Tactics  []*Tactic   `json:"tactics"`
	Training []*Training `json:"training"`
}

func (b *Storage) NewUserActionsByVerse(verse int) (*UserActions, error) {
	var err error
	var userActions UserActions
	if userActions.Tactics, err = b.TacticsByVerse(verse); err != nil {
		return nil, err
	}
	return &userActions, nil
}

func (b *UserActions) Marshal() ([]byte, error) {
	if b.Tactics == nil {
		b.Tactics = make([]*Tactic, 0)
	}
	if b.Training == nil {
		b.Training = make([]*Training, 0)
	}
	return json.Marshal(b)
}

func (b *UserActions) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &b)
}
