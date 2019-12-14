package useractions

import (
	"archive/tar"
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"

	"github.com/freeverseio/crypto-soccer/go/relay/storage"

	shell "github.com/ipfs/go-ipfs-api"
)

type UserActions struct {
	Tactics   []*storage.Tactic   `json:"tactics"`
	Trainings []*storage.Training `json:"trainings"`
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
	if b.Trainings == nil {
		b.Trainings = make([]*storage.Training, 0)
	}
	return json.Marshal(b)
}

func (b *UserActions) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &b)
}

func (b *UserActions) IpfsPush() (string, error) {
	sh := shell.NewShell("localhost:5001")
	buf, err := b.Marshal()
	if err != nil {
		return "", err
	}
	return sh.Add(bytes.NewReader(buf), shell.Pin(true))
}

func (b *UserActions) IpfsPull(cid string) error {
	sh := shell.NewShell("localhost:5001")
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
