package process_test

import (
	//	"fmt"

	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/process"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/testutils"
)

// TODO commented cause I dunno what's happening
// func TestSyncTeamWithNoTeam(t *testing.T) {
// 	storage, err := storage.NewSqlite3("../sql/00_schema.sql")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	blockchain := testutils.DefaultSimulatedBlockchain()

// 	p := NewEventProcessor(nil, storage, blockchain.Assets, blockchain.States, blockchain.Leagues)
// 	p.Process()

// 	count, err := storage.TeamCount()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if count != 0 {
// 		t.Fatalf("Expected 0 received %v", count)
// 	}
// }

func TestSyncTeams(t *testing.T) {
	storage, err := storage.NewSqlite3("../sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	ganache := testutils.NewGanache()
	_ = storage

	owner := ganache.CreateAccountWithBalance("1000000000000000000") // 1 eth
	ganache.DeployContracts(owner)

	alice := ganache.CreateAccountWithBalance("50000000000000000000") // 50 eth
	bob := ganache.CreateAccountWithBalance("50000000000000000000")   // 50 eth
	carol := ganache.CreateAccountWithBalance("50000000000000000000") // 50 eth

	timezoneIdx := uint8(1)
	countryIdx := big.NewInt(0)
	firstTeamID, err := ganache.Leagues.EncodeTZCountryAndVal(
		nil,
		timezoneIdx,
		countryIdx,
		big.NewInt(0),
	)
	if err != nil {
		t.Fatal(err)
	}
	if firstTeamID.String() != "274877906944" {
		t.Fatalf("Expected 274877906944 but received %v", firstTeamID.String())
	}

	_ = alice
	_ = bob
	_ = carol
	//ganache.CreateTeam("A", alice)
	//ganache.CreateTeam("B", bob)
	//ganache.CreateTeam("C", carol)

	p := process.NewGanacheEventProcessor(ganache.Client, storage, ganache.Engine, ganache.Leagues, ganache.Updates)

	if err := p.Process(); err != nil {
		t.Fatal(err)
	} else {
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		p.Process()
		if count, err := storage.TimezoneCount(); err != nil {
			t.Fatal(err)
		} else if count != 1 {
			t.Fatalf("Expected 1 time zones at time of creation,  actual %v", count)
		}

		if count, err := storage.CountryCount(); err != nil {
			t.Fatal(err)
		} else if count != 1 {
			t.Fatalf("Expected 1 countries at time of creation,  actual %v", count)
		}

		if count, err := storage.TeamCount(); err != nil {
			t.Fatal(err)
		} else if count != 128 {
			t.Fatalf("Expected 128 actual %v", count)
		}
		if count, err := storage.PlayerCount(); err != nil {
			t.Fatal(err)
		} else if count != 128*18 {
			t.Fatalf("Expected 128*18=2304 actual %v", count)
		}
	}

	_, err = ganache.Leagues.TransferFirstBotToAddr(
		bind.NewKeyedTransactor(owner),
		timezoneIdx,
		countryIdx,
		ganache.Public(alice),
	)
	if err != nil {
		t.Fatal(err)
	}

	teamOwner, err := ganache.Leagues.GetOwnerTeam(nil, firstTeamID)
	if err != nil {
		t.Fatal(err)
	}
	if teamOwner != ganache.Public(alice) {
		t.Fatalf("first team owner is %v", teamOwner.String())
	}

	err = p.Process()
	if err != nil {
		t.Fatal(err)
	}

	team, err := storage.GetTeam(firstTeamID)
	if err != nil {
		t.Fatal(err)
	}

	if team.State.Owner != ganache.Public(alice).String() {
		t.Fatalf("database owner of the first team is %v", team.State.Owner)
	}

	// play
	for i := 0; i < 24*4; i++ {
		var root [32]byte
		_, err = ganache.Updates.SubmitActionsRoot(
			bind.NewKeyedTransactor(owner),
			root,
		)
		if err != nil {
			t.Fatal(err)
		}
	}

	p.Process()

	//fmt.Println("owner: ", ganache.Public(ganache.Owner).Hex())
	//fmt.Println("alice: ", ganache.Public(alice).Hex())
	//fmt.Println("bob: ", ganache.Public(bob).Hex())
	//fmt.Println("carol: ", ganache.Public(carol).Hex())
	//// tema A
	//if team, err := storage.GetTeam(1); err != nil {
	//	t.Fatal(err)
	//} else if team.Id != 1 {
	//	t.Fatalf("Expected 1 result %v", team.Id)
	//} else if team.Name != "A" {
	//	t.Fatalf("Expected A result %v", team.Name)
	//} else if state, err := storage.GetTeam(1); err != nil {
	//	t.Fatal(err)
	//} else if state.State.Owner != ganache.Public(alice).Hex() {
	//	t.Fatalf("Expecting team A to belong to Alice %v : %v", state.State.Owner, ganache.Public(alice).Hex())
	//}
	//// team B
	//if team, err := storage.GetTeam(2); err != nil {
	//	t.Fatal(err)
	//} else if team.Id != 2 {
	//	t.Fatalf("Expected 2 result %v", team.Id)
	//} else if team.Name != "B" {
	//	t.Fatalf("Expected B result %v", team.Name)
	//} else if state, err := storage.GetTeam(2); err != nil {
	//	t.Fatal(err)
	//} else if state.State.Owner != ganache.Public(bob).Hex() {
	//	t.Fatalf("Expecting team B to belong to Bob %v : %v", state.State.Owner, ganache.Public(bob).Hex())
	//}
	//// team C
	//if team, err := storage.GetTeam(3); err != nil {
	//	t.Fatal(err)
	//} else if team.Id != 3 {
	//	t.Fatalf("Expected 3 result %v", team.Id)
	//} else if team.Name != "C" {
	//	t.Fatalf("Expected C result %v", team.Name)
	//} else if state, err := storage.GetTeam(3); err != nil {
	//	t.Fatal(err)
	//} else if state.State.Owner != ganache.Public(carol).Hex() {
	//	t.Fatalf("Expecting team A to belong to Carol %v : %v", state.State.Owner, ganache.Public(carol).Hex())
	//}

	//if count, err := storage.PlayerCount(); err != nil {
	//	t.Fatal(err)
	//} else if count != 54 {
	//	t.Fatalf("Expected 54 players actual %v", count)
	//} else {
	//	for i := 1; i <= 18; i++ {
	//		if result, err := storage.GetPlayer(uint64(i)); err != nil {
	//			t.Fatal(err)
	//		} else if result.State.TeamId == uint64(0) || result.Id != uint64(i) {
	//			t.Fatalf("Expecting player %v state to be non 0 actual %v", i, result)
	//		}
	//	}
	//}

	//ganache.CreateTeam("D", alice)
	//if err := p.Process(); err != nil {
	//	t.Fatal(err)
	//} else {
	//	if count, err := storage.TeamCount(); err != nil {
	//		t.Fatal(err)
	//	} else if count != 4 {
	//		t.Fatalf("Expected 4 actual %v", count)
	//	}
	//}
	//if team, err := storage.GetTeam(4); err != nil {
	//	t.Fatal(err)
	//} else if team.Id != 4 {
	//	t.Fatalf("Expected 4 result %v", team.Id)
	//} else if team.Name != "D" {
	//	t.Fatalf("Expected D result %v", team.Name)
	//}

	//if count, err := storage.PlayerCount(); err != nil {
	//	t.Fatal(err)
	//} else if count != 72 {
	//	t.Fatalf("Expected 72 players actual %v", count)
	//}

	//// ganache.CreateLeague([]int64{1, 2, 3, 4}, alice)
	//ganache.Advance(3) // advance 3 blocks
	p.Process()
}
