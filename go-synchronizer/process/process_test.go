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

	p := NewEventProcessor(nil, storage, blockchain.Assets, blockchain.States, blockchain.Leagues)
	p.Process()

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

	p := NewEventProcessor(ganache.Client, storage, ganache.Assets, ganache.States, ganache.Leagues)

	if err := p.Process(); err != nil {
		t.Fatal(err)
	} else {
		if count, err := storage.TeamCount(); err != nil {
			t.Fatal(err)
		} else if count != 3 {
			t.Fatalf("Expected 3 actual %v", count)
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

	if count, err := storage.PlayerCount(); err != nil {
		t.Fatal(err)
	} else if count != 33 {
		t.Fatalf("Expected 33 players actual %v", count)
	} else {
		for i := 1; i <= 33; i++ {
			if result, err := storage.GetPlayer(uint64(i)); err != nil {
				t.Fatal(err)
			} else if result.State == "0" || result.Id != uint64(i) {
				t.Fatalf("Expecting player %v state to be non 0 actual %v", i, result)
			}
		}
	}

	ganache.CreateTeam("D", alice)
	if err := p.Process(); err != nil {
		t.Fatal(err)
	} else {
		if count, err := storage.TeamCount(); err != nil {
			t.Fatal(err)
		} else if count != 4 {
			t.Fatalf("Expected 4 actual %v", count)
		}
	}
	if team, err := storage.GetTeam(4); err != nil {
		t.Fatal(err)
	} else if team.Id != 4 {
		t.Fatalf("Expected 4 result %v", team.Id)
	} else if team.Name != "D" {
		t.Fatalf("Expected D result %v", team.Name)
	}

	if count, err := storage.PlayerCount(); err != nil {
		t.Fatal(err)
	} else if count != 44 {
		t.Fatalf("Expected 44 players actual %v", count)
	}

	ganache.CreateLeague([]int64{1, 2, 3, 4}, alice)
}
