package ipfscluster_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/useractions"

	"github.com/freeverseio/crypto-soccer/go/useractions/ipfscluster"
)

func TestUserActionsPublishService(t *testing.T) {
	service := ipfscluster.NewUserActionsPublishService("/ip4/127.0.0.1/tcp/5001")
	useractions.TestUserActionsPublishService(t, service)
}
