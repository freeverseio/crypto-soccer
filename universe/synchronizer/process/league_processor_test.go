package process_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/freeverseio/crypto-soccer/universe/synchronizer/contracts/leagues"

	"github.com/freeverseio/crypto-soccer/universe/synchronizer/contracts/updates"
	"github.com/freeverseio/crypto-soccer/universe/synchronizer/process"
	"github.com/freeverseio/crypto-soccer/universe/synchronizer/storage"
	"github.com/freeverseio/crypto-soccer/universe/synchronizer/testutils"
)

func TestCreateMAtchSeed(t *testing.T) {
	ganache, err := testutils.NewBlockchainNode()
	if err != nil {
		t.Fatal(err)
	}
	ganache.DeployContracts(ganache.Owner)
	processor, err := process.NewLeagueProcessor(ganache.Engine, ganache.Leagues, nil)
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
func TestProcessInvalidTimezone(t *testing.T) {
	sto, err := storage.NewSqlite3("../sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	ganache, err := testutils.NewBlockchainNode()
	if err != nil {
		t.Fatal(err)
	}
	ganache.DeployContracts(ganache.Owner)
	processor, err := process.NewLeagueProcessor(ganache.Engine, ganache.Leagues, sto)
	if err != nil {
		t.Fatal(err)
	}
	return // TODO check
	var event updates.UpdatesActionsSubmission
	event.TimeZone = 25
	err = processor.Process(event)
	if err == nil {
		t.Fatal("processing invalid timezone")
	}
}

func TestProcess(t *testing.T) {
	sto, err := storage.NewSqlite3("../sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	bc, err := testutils.NewBlockchainNode()
	if err != nil {
		t.Fatal(err)
	}
	err = bc.DeployContracts(bc.Owner)
	if err != nil {
		t.Fatal(err)
	}
	timezoneIdx := uint8(1)
	err = bc.InitOneTimezone(timezoneIdx)
	if err != nil {
		t.Fatal(err)
	}
	divisionCreationProcessor, err := process.NewDivisionCreationProcessor(sto, bc.Leagues)
	if err != nil {
		t.Fatal(err)
	}
	countryIdx := big.NewInt(0)
	divisionIdx := big.NewInt(0)
	err = divisionCreationProcessor.Process(leagues.LeaguesDivisionCreation{timezoneIdx, countryIdx, divisionIdx, types.Log{}})
	if err != nil {
		t.Fatal(err)
	}

	processor, err := process.NewLeagueProcessor(bc.Engine, bc.Leagues, sto)
	if err != nil {
		t.Fatal(err)
	}
	var event updates.UpdatesActionsSubmission
	event.Day = 1
	event.TimeZone = timezoneIdx
	event.TurnInDay = 1
	day := uint8(1)
	turnInDay := uint8(1)
	seed := [32]byte{}
	err = processor.Process(updates.UpdatesActionsSubmission{timezoneIdx, day, turnInDay, seed, nil, types.Log{}})
	if err != nil {
		t.Fatal(err)
	}
}
