package process_test

import (
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/names"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/process"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"
	"github.com/freeverseio/crypto-soccer/go/testutils"
)

func TestCreateMatchSeed(t *testing.T) {
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
	tx, err := universedb.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

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
		namesdb,
	)
	if err != nil {
		t.Fatal(err)
	}
	teamState, err := processor.GetTeamState(tx, big.NewInt(3))
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

func TestProcessMatchEvents(t *testing.T) {
	tx, err := universedb.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
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
		namesdb,
	)
	if err != nil {
		t.Fatal(err)
	}
	match := storage.Match{}
	match.HomeTeamID = big.NewInt(274877906944)
	match.VisitorTeamID = big.NewInt(274877906945)
	match.HomeMatchLog = big.NewInt(0)
	match.VisitorMatchLog = big.NewInt(0)
	states := [2][25]*big.Int{}
	states[0], err = processor.GetTeamState(tx, match.HomeTeamID)
	if err != nil {
		t.Fatal(err)
	}
	states[1], err = processor.GetTeamState(tx, match.VisitorTeamID)
	if err != nil {
		t.Fatal(err)
	}
	tactics, err := processor.GetMatchTactics(match.HomeTeamID, match.VisitorTeamID)
	if err != nil {
		t.Fatal(err)
	}
	seed := [32]byte{0x0}
	matchSeed, err := processor.GenerateMatchSeed(seed, match.HomeTeamID, match.VisitorTeamID)
	if err != nil {
		t.Fatal(err)
	}
	startTime := big.NewInt(555)
	is2ndHalf := false
	events, err := processor.ProcessMatchEvents(
		match,
		states,
		tactics,
		matchSeed,
		startTime,
		is2ndHalf,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 16 {
		t.Fatalf("Wrong length of events  in 1st half %v", len(events))
	}
	is2ndHalf = true
	events, err = processor.ProcessMatchEvents(
		match,
		states,
		tactics,
		matchSeed,
		startTime,
		is2ndHalf,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 16 {
		t.Fatalf("Wrong length of events in 2nd half %v", len(events))
	}
}
