package memory

import (
	"encoding/hex"

	"github.com/freeverseio/crypto-soccer/go/useractions"
)

type UserActionsPublishService struct {
	data map[[32]byte]useractions.UserActions
}

func NewUserActionsPublishService() *UserActionsPublishService {
	return &UserActionsPublishService{
		data: make(map[[32]byte]useractions.UserActions),
	}

}

func (b UserActionsPublishService) Publish(actions useractions.UserActions) (string, error) {
	root, err := actions.Root()
	if err != nil {
		return "", err
	}
	b.data[root] = actions
	return hex.EncodeToString(root[:]), nil
}

func (b UserActionsPublishService) Retrive(id string) (*useractions.UserActions, error) {
	root, err := hex.DecodeString(id)
	if err != nil {
		return nil, err
	}
	var key [32]byte
	copy(key[:], root)
	result := b.data[key]
	return &result, nil
}
