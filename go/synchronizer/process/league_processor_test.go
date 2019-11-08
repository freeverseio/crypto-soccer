package process_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/freeverseio/crypto-soccer/go/contracts/assets"
	"github.com/freeverseio/crypto-soccer/go/contracts/updates"
	relay "github.com/freeverseio/crypto-soccer/go/relay/storage"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/process"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"
	"github.com/freeverseio/crypto-soccer/go/testutils"
)

func TestCreateMAtchSeed(t *testing.T) {
	ganache, err := testutils.NewBlockchainNode()
	if err != nil {
		t.Fatal(err)
	}
	ganache.DeployContracts(ganache.Owner)
	processor, err := process.NewLeagueProcessor(
		ganache.Engine,
		ganache.EnginePreComp,
		ganache.Assets,
		ganache.Leagues,
		ganache.Evolution,
		nil,
		nil,
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
func TestProcessInvalidTimezone(t *testing.T) {
	universedb, err := storage.NewSqlite3("../../../universe.db/00_schema.sql")
	relaydb, err := relay.NewSqlite3("../../../relay.db/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	ganache, err := testutils.NewBlockchainNode()
	if err != nil {
		t.Fatal(err)
	}
	ganache.DeployContracts(ganache.Owner)
	processor, err := process.NewLeagueProcessor(
		ganache.Engine,
		ganache.EnginePreComp,
		ganache.Assets,
		ganache.Leagues,
		ganache.Evolution,
		universedb,
		relaydb,
	)
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

func TestPlayHalfMatch(t *testing.T) {
	bc, err := testutils.NewBlockchainNode()
	if err != nil {
		t.Fatal(err)
	}
	err = bc.DeployContracts(bc.Owner)
	if err != nil {
		t.Fatal(err)
	}
	seed, _ := new(big.Int).SetString("79914733518162338039108289434138268582051932490617416751079415725906677907530", 10)
	matchStartTime := big.NewInt(1570147200)
	matchLog := [2]*big.Int{big.NewInt(0), big.NewInt(0)}
	matchBools := [3]bool{false, false, false}
	var states [2][25]*big.Int
	for i := 0; i < 2; i++ {
		for j := 0; j < 25; j++ {
			states[i][j], _ = new(big.Int).SetString("713624055286353394965726120199142814938406092850", 10)
		}
	}
	var tactics [2]*big.Int
	tactics[0], _ = new(big.Int).SetString("117587469164573768163156115324928", 10)
	tactics[1], _ = new(big.Int).SetString("117587469164573768163156115324928", 10)
	logs, err := bc.Engine.PlayHalfMatch(
		&bind.CallOpts{},
		seed,
		matchStartTime,
		states,
		tactics,
		matchLog,
		matchBools,
	)
	if err != nil {
		t.Fatal(err)
	}
	if logs[0].String() != "1645504557572631091828279073381766814124583466071029250581856256" {
		t.Fatalf("Received %v", logs[0].String())
	}
}

func TestLeagueProcessMatch(t *testing.T) {
	universedb, err := storage.NewSqlite3("../../../universe.db/00_schema.sql")
	relaydb, err := relay.NewSqlite3("../../../relay.db/00_schema.sql")
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
	divisionCreationProcessor, err := process.NewDivisionCreationProcessor(universedb, bc.Assets, bc.Leagues)
	if err != nil {
		t.Fatal(err)
	}
	countryIdx := big.NewInt(0)
	divisionIdx := big.NewInt(0)
	err = divisionCreationProcessor.Process(assets.AssetsDivisionCreation{timezoneIdx, countryIdx, divisionIdx, types.Log{}})
	if err != nil {
		t.Fatal(err)
	}

	processor, err := process.NewLeagueProcessor(
		bc.Engine,
		bc.EnginePreComp,
		bc.Assets,
		bc.Leagues,
		bc.Evolution,
		universedb,
		relaydb,
	)
	if err != nil {
		t.Fatal(err)
	}
	day := uint8(0)
	turnInDay := uint8(0)
	seed := [32]byte{}
	err = processor.Process(updates.UpdatesActionsSubmission{
		timezoneIdx,
		day,
		turnInDay,
		seed,
		big.NewInt(10),
		types.Log{},
	})
	if err != nil {
		t.Fatal(err)
	}
	turnInDay = 1
	err = processor.Process(updates.UpdatesActionsSubmission{
		timezoneIdx,
		day,
		turnInDay,
		seed,
		big.NewInt(10),
		types.Log{},
	})
	if err != nil {
		t.Fatal(err)
	}
}
