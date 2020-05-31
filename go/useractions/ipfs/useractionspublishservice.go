package ipfs

import (
	"archive/tar"
	"bytes"
	"context"
	"io/ioutil"

	"github.com/freeverseio/crypto-soccer/go/useractions"
	shell "github.com/ipfs/go-ipfs-api"
)

type UserActionsPublishService struct {
	url string
}

func NewUserActionsPublishService(endpoint string) *UserActionsPublishService {
	return &UserActionsPublishService{
		url: endpoint,
	}
}

func (b UserActionsPublishService) Publish(actions useractions.UserActions) (string, error) {
	sh := shell.NewShell(b.url)
	buf, err := actions.Marshal()
	if err != nil {
		return "", err
	}
	return sh.Add(bytes.NewReader(buf), shell.Pin(true))
}

func (b UserActionsPublishService) Retrive(cid string) (*useractions.UserActions, error) {
	var ua useractions.UserActions
	sh := shell.NewShell(b.url)
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
