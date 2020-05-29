package useractions

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/freeverseio/crypto-soccer/go/storage"

	shell "github.com/ipfs/go-ipfs-api"
	cluster "github.com/ipfs/ipfs-cluster/api/rest/client"
	ma "github.com/multiformats/go-multiaddr"
)

type UserActions struct {
	Tactics   []storage.Tactic   `json:"tactics"`
	Trainings []storage.Training `json:"trainings"`
}

type UserActionsPublishService interface {
	Publish(actions UserActions) (string, error)
	Retrive(cid string) (*UserActions, error)
}

func (b *UserActions) Equal(actions *UserActions) bool {
	if len(b.Tactics) != len(actions.Tactics) {
		return false
	}
	if len(b.Trainings) != len(actions.Trainings) {
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
	return true
}

func New() *UserActions {
	return &UserActions{}
}

func newShell(url string) (*shell.Shell, error) {
	maddr, err := ma.NewMultiaddr(url)
	if err != nil {
		return nil, err
	}
	cfg := &cluster.Config{
		APIAddr: maddr,
	}
	cfg.SSL = false
	cfg.NoVerifyCert = true
	client, err := cluster.NewDefaultClient(cfg)
	if err != nil {
		return nil, err
	}
	return client.IPFS(context.Background()), nil
}

func NewFromStorage(tx *sql.Tx, timezone int) (*UserActions, error) {
	var err error
	var ua UserActions
	if ua.Tactics, err = storage.TacticsByTimezone(tx, timezone); err != nil {
		return nil, err
	}
	if ua.Trainings, err = storage.TrainingsByTimezone(tx, timezone); err != nil {
		return nil, err
	}
	return &ua, nil
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
	return json.Marshal(b)
}

func (b *UserActions) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &b)
}
