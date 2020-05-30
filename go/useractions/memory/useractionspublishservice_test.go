package memory_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/useractions/memory"
	"github.com/freeverseio/crypto-soccer/go/useractions/useractionstest"
)

func TestUserActionsPublishService(t *testing.T) {
	service := memory.NewUserActionsPublishService()
	useractionstest.TestUserActionsPublishService(t, service)
}
