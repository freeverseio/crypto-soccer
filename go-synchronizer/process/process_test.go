package process

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/testutils"
)

func TestSyncTeamWithNoTeam(t *testing.T) {
	storage, err := storage.NewSqlite3("../sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	blockchain := testutils.DefaultSimulatedBlockchain()

	Process(blockchain.Assets, storage, nil)

	count, err := storage.TeamCount()
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("Expected 0 received %v", count)
	}
}

func TestSyncTeams(t *testing.T) {
	storage, err := storage.NewSqlite3("../sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	ganache := testutils.NewGanache()

	owner := ganache.CreateAccountWithBalance("1000000000000000000") // 1 eth
	ganache.DeployContracts(owner)

	alice := ganache.CreateAccountWithBalance("50000000000000000000") // 50 eth
	bob := ganache.CreateAccountWithBalance("50000000000000000000")   // 50 eth
	carol := ganache.CreateAccountWithBalance("50000000000000000000") // 50 eth

	ganache.CreateTeam("A", alice)
	ganache.CreateTeam("B", bob)
	ganache.CreateTeam("C", carol)

	if err := Process(ganache.Assets, storage, ganache.Client); err != nil {
		t.Fatal(err)
	} else {
		if count, err := storage.TeamCount(); err != nil {
			t.Fatal(err)
		} else if count != 3 {
			t.Fatalf("Expected 3 received %v", count)
		}
	}

	if team, err := storage.GetTeam(1); err != nil {
		t.Fatal(err)
	} else if team.Id != 1 {
		t.Fatalf("Expected 1 result %v", team.Id)
	} else if team.Name != "A" {
		t.Fatalf("Expected A result %v", team.Name)
	}
	if team, err := storage.GetTeam(2); err != nil {
		t.Fatal(err)
	} else if team.Id != 2 {
		t.Fatalf("Expected 2 result %v", team.Id)
	} else if team.Name != "B" {
		t.Fatalf("Expected B result %v", team.Name)
	}
	if team, err := storage.GetTeam(3); err != nil {
		t.Fatal(err)
	} else if team.Id != 3 {
		t.Fatalf("Expected 3 result %v", team.Id)
	} else if team.Name != "C" {
		t.Fatalf("Expected C result %v", team.Name)
	}

	ganache.CreateTeam("D", alice)
	if err := Process(ganache.Assets, storage, ganache.Client); err != nil {
		t.Fatal(err)
	} else {
		if count, err := storage.TeamCount(); err != nil {
			t.Fatal(err)
		} else if count != 4 {
			t.Fatalf("Expected 4 received %v", count)
		}
	}
	if team, err := storage.GetTeam(4); err != nil {
		t.Fatal(err)
	} else if team.Id != 4 {
		t.Fatalf("Expected 4 result %v", team.Id)
	} else if team.Name != "D" {
		t.Fatalf("Expected D result %v", team.Name)
	}
}
