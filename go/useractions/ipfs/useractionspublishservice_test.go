package ipfs_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/useractions"
	"github.com/freeverseio/crypto-soccer/go/useractions/ipfs"
)

func TestUserActionsPublishService(t *testing.T) {
	service := ipfs.NewUserActionsPublishService("/ip4/127.0.0.1/tcp/5001")
	useractions.TestUserActionsPublishService(t, service)
}
