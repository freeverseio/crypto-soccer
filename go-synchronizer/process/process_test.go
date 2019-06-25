package process

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage/memory"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/testutils"
)

func TestSyncTeam(t *testing.T) {
	storage := memory.New()
	blockchain := testutils.CryptosoccerNew(t)

	Process(blockchain.AssetsContract, storage)

	count, err := storage.TeamCount()
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("Expected 0 received %v", count)
	}
}
