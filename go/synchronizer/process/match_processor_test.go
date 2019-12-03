package process_test

import (
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/names"
	relay "github.com/freeverseio/crypto-soccer/go/relay/storage"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/process"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"
	"github.com/freeverseio/crypto-soccer/go/testutils"
)

func TestCreateMatchSeed(t *testing.T) {
	universedb, err := storage.NewSqlite3("../../../universe.db/00_schema.sql")
	relaydb, err := relay.NewSqlite3("../../../universe.db/00_schema.sql")
	namesdb, err := names.New("../../names/sql/names.db")
	if err != nil {
		t.Fatal(err)
	}
	bc, err := testutils.NewBlockchainNode()
	if err != nil {
		t.Fatal(err)
	}
	bc.DeployContracts(bc.Owner)
	processor, err := process.NewMatchProcessor(
		bc.Contracts,
		universedb,
		relaydb,
		namesdb,
	)
	if err != nil {
		t.Fatal(err)
	}
	seed := [32]byte{0x0}
	homeTeamID := big.NewInt(3)
	visitorTeamID := big.NewInt(5)
	result, err := processor.GenerateMatchSeed(seed, homeTeamID, visitorTeamID)
	if err != nil {
		t.Fatal(err)
	}
	if result.Text(16) != "33c646d693b716acb3a01ae35dd9ed16191786670a88f4c086b7223851a750d" {
		t.Fatalf("Received %v", result.Text(16))
	}
}

func TestGetPlayerState(t *testing.T) {
	universedb, err := storage.NewSqlite3("../../../universe.db/00_schema.sql")
	relaydb, err := relay.NewSqlite3("../../../universe.db/00_schema.sql")
	namesdb, err := names.New("../../names/sql/names.db")
	if err != nil {
		t.Fatal(err)
	}
	bc, err := testutils.NewBlockchainNode()
	if err != nil {
		t.Fatal(err)
	}
	bc.DeployContracts(bc.Owner)
	processor, err := process.NewMatchProcessor(
		bc.Contracts,
		universedb,
		relaydb,
		namesdb,
	)
	if err != nil {
		t.Fatal(err)
	}
	teamState, err := processor.GetTeamState(big.NewInt(3))
	if err != nil {
		t.Fatal(err)
	}
	if len(teamState) != 25 {
		t.Fatalf("Wrong team state count %v", len(teamState))
	}
	for _, state := range teamState {
		if state.String() != "0" {
			t.Fatalf("Wrong state %v", state)
		}
	}
}
