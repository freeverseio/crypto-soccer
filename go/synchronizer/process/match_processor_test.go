package process_test

import (
	"math/big"
	"testing"

	relay "github.com/freeverseio/crypto-soccer/go/relay/storage"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/process"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"
	"github.com/freeverseio/crypto-soccer/go/testutils"
)

func TestCreateMatchSeed(t *testing.T) {
	universedb, err := storage.NewSqlite3("../../../universe.db/00_schema.sql")
	relaydb, err := relay.NewSqlite3("../../../relay.db/00_schema.sql")
	ganache, err := testutils.NewBlockchainNode()
	if err != nil {
		t.Fatal(err)
	}
	ganache.DeployContracts(ganache.Owner)
	processor, err := process.NewMatchProcessor(
		universedb,
		relaydb,
		ganache.Assets,
		ganache.Leagues,
		ganache.Evolution,
		ganache.Engine,
		ganache.EnginePreComp,
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
