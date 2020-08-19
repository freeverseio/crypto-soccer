package useractions

import (
	"crypto/sha256"
	"encoding/json"
	"errors"

	"github.com/freeverseio/crypto-soccer/go/storage"
	"github.com/freeverseio/crypto-soccer/go/useractions/orgmap"
)

type UserActions struct {
	Tactics        []storage.Tactic        `json:"tactics"`
	Trainings      []storage.Training      `json:"trainings"`
	OrgMapDenyList []orgmap.OrgMapDenyList `json:"orgmapdenylist"`
}

type UserActionsPublishService interface {
	Publish(actions UserActions) (string, error)
	Retrive(cid string) (*UserActions, error)
}

type UserActionsStorageService interface {
	UserActionsByTimezone(timezone int) (*UserActions, error)
	Insert(actions UserActions) error
	InsertHistory(blockNumber uint64, actions UserActions) error
}

func (b *UserActions) Hash() ([32]byte, error) {
	var result [32]byte
	h := sha256.New()
	buf, err := b.Marshal()
	if err != nil {
		return result, err
	}
	h.Write(buf)
	hash := h.Sum(nil)
	if len(hash) != 32 {
		return result, errors.New("Hash is not 32 byte")
	}
	copy(result[:], hash)
	return result, nil
}

func (b UserActions) Root() ([32]byte, error) {
	return b.Hash()
}

func (b *UserActions) Marshal() ([]byte, error) {
	if b.Tactics == nil {
		b.Tactics = make([]storage.Tactic, 0)
	}
	if b.Trainings == nil {
		b.Trainings = make([]storage.Training, 0)
	}
	if b.OrgMapDenyList == nil {
		b.OrgMapDenyList = make([]orgmap.OrgMapDenyList, 0)
	}
	return json.Marshal(b)
}

func (b *UserActions) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &b)
}

func (b *UserActions) Equal(actions *UserActions) bool {
	if len(b.Tactics) != len(actions.Tactics) {
		return false
	}
	if len(b.Trainings) != len(actions.Trainings) {
		return false
	}
	if len(b.OrgMapDenyList) != len(actions.OrgMapDenyList) {
		return false
	}

	for i := range b.Tactics {
		if b.Tactics[i] != actions.Tactics[i] {
			return false
		}
	}
	for i := range b.Trainings {
		if b.Trainings[i] != actions.Trainings[i] {
			return false
		}
	}
	for i := range b.OrgMapDenyList {
		if b.OrgMapDenyList[i] != actions.OrgMapDenyList[i] {
			return false
		}
	}
	return true
}
