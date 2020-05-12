package useractions

import (
	"archive/tar"
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/freeverseio/crypto-soccer/go/storage"

	shell "github.com/ipfs/go-ipfs-api"
	cluster "github.com/ipfs/ipfs-cluster/api/rest/client"
	ma "github.com/multiformats/go-multiaddr"
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

func New() *UserActions {
	return &UserActions{}
}

func newShell(url string) *shell.Shell {
	useIpfsCluster := false
	if !useIpfsCluster {
		return shell.NewShell(url)
	}
	if len(url) == 0 {
		url = "/ip4/127.0.0.1/tcp/5001" // localhost
	}
	maddr, err := ma.NewMultiaddr(url)
	if err != nil {
		panic(err)
	}
	cfg := &cluster.Config{
		APIAddr: maddr,
		//Username: *username,
		//Password: *pw,
	}
	cfg.SSL = false //true
	cfg.NoVerifyCert = true
	client, err := cluster.NewDefaultClient(cfg)
	if err != nil {
		panic(err)
	}
	return client.IPFS(context.Background())
}

func NewFromIpfs(url string, cid string) (*UserActions, error) {
	var ua UserActions
	sh := newShell(url)
	resp, err := sh.Request("get", cid).Option("create", true).Send(context.Background())
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	if resp.Error != nil {
		return nil, resp.Error
	}

	tr := tar.NewReader(resp.Output)
	_, err = tr.Next()
	if err != nil {
		return nil, err
	}
	buf, err := ioutil.ReadAll(tr)
	if err != nil {
		return nil, err
	}
	if err = ua.Unmarshal(buf); err != nil {
		return nil, err
	}
	return &ua, nil
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

func (b *UserActions) ToIpfs(url string) (string, error) {
	sh := newShell(url)
	buf, err := b.Marshal()
	if err != nil {
		return "", err
	}
	return sh.Add(bytes.NewReader(buf), shell.Pin(true))
}
