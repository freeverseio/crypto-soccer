package ipfs_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/useractions/ipfs"
	"github.com/freeverseio/crypto-soccer/go/useractions/useractionstest"
)

func TestUserActionsPublishService(t *testing.T) {
	service := ipfs.NewUserActionsPublishService("/ip4/127.0.0.1/tcp/5001")
	useractionstest.TestUserActionsPublishService(t, service)
}
