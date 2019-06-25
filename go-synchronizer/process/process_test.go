package process

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage/memory"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/testutils"
)

func TestSyncTeamWithNoTeam(t *testing.T) {
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

func TestSyncTeams(t *testing.T) {
	storage := memory.New()
	blockchain := testutils.CryptosoccerNew(t)
	blockchain.AssetsContract.CreateTeam(blockchain.Opts, "Barca", blockchain.Opts.From)
	blockchain.AssetsContract.CreateTeam(blockchain.Opts, "Madrid", blockchain.Opts.From)
	blockchain.AssetsContract.CreateTeam(blockchain.Opts, "Venezia", blockchain.Opts.From)
	blockchain.Commit()

	err := Process(blockchain.AssetsContract, storage)
	if err != nil {
		t.Fatal(err)
	}

	count, err := storage.TeamCount()
	if err != nil {
		t.Fatal(err)
	}
	if count != 3 {
		t.Fatalf("Expected 3 received %v", count)
	}
	team, err := storage.GetTeam(1)
	if err != nil {
		t.Fatal(err)
	}
	if team.Id != 1 {
		t.Fatalf("Expected 1 result %v", team.Id)
	}
	if team.Name != "Barca" {
		t.Fatalf("Expected Barca result %v", team.Name)
	}
	team, err = storage.GetTeam(2)
	if err != nil {
		t.Fatal(err)
	}
	if team.Id != 2 {
		t.Fatalf("Expected 2 result %v", team.Id)
	}
	if team.Name != "Madrid" {
		t.Fatalf("Expected Madrid result %v", team.Name)
	}
	team, err = storage.GetTeam(3)
	if err != nil {
		t.Fatal(err)
	}
	if team.Id != 3 {
		t.Fatalf("Expected 3 result %v", team.Id)
	}
	if team.Name != "Venezia" {
		t.Fatalf("Expected Venezia result %v", team.Name)
	}
}
