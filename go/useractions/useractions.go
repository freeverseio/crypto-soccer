package useractions

import (
	"archive/tar"
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"io/ioutil"

	"github.com/freeverseio/crypto-soccer/go/relay/storage"

	shell "github.com/ipfs/go-ipfs-api"
)

type UserActions struct {
	Tactics   []storage.Tactic   `json:"tactics"`
	Trainings []storage.Training `json:"trainings"`
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

func (b *UserActions) PullFromStorage(tx *sql.Tx, verse uint64) error {
	var err error
	if b.Tactics, err = storage.TacticsByVerse(tx, verse); err != nil {
		return err
	}
	if b.Trainings, err = storage.TrainingByVerse(tx, verse); err != nil {
		return err
	}
	return nil
}

func (b *UserActions) Hash() ([]byte, error) {
	h := sha256.New()
	buf, err := b.Marshal()
	if err != nil {
		return nil, err
	}
	h.Write(buf)
	return h.Sum(nil), nil
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

func (b *UserActions) PushToIpfs(url string) (string, error) {
	sh := shell.NewShell(url)
	buf, err := b.Marshal()
	if err != nil {
		return "", err
	}
	return sh.Add(bytes.NewReader(buf), shell.Pin(true))
}

func (b *UserActions) PullFromIpfs(url string, cid string) error {
	sh := shell.NewShell(url)
	resp, err := sh.Request("get", cid).Option("create", true).Send(context.Background())
	if err != nil {
		return err
	}
	defer resp.Close()

	if resp.Error != nil {
		return resp.Error
	}

	tr := tar.NewReader(resp.Output)
	_, err = tr.Next()
	if err != nil {
		return err
	}
	buf, err := ioutil.ReadAll(tr)
	if err != nil {
		return err
	}
	return b.Unmarshal(buf)
}
