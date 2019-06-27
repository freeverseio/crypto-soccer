package scanners

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go-synchronizer/testutils"
)

func TestEmptyContract(t *testing.T) {
	blockchain := testutils.DefaultSimulatedBlockchain()
	events, err := ScanTeamCreated(blockchain.Assets, nil)
	testutils.AssertNoErr(err)
	if len(events) != 0 {
		t.Fatalf("Scanning empty Assets contract returned %v events", len(events))
	}
}

func TestCreateTeam(t *testing.T) {
	blockchain := testutils.DefaultSimulatedBlockchain()
	alice := blockchain.CreateAccountWithBalance("1000000000000000000") // 1 eth
	blockchain.CreateTeam("Barca", alice)

	events, err := ScanTeamCreated(blockchain.Assets, nil)
	testutils.AssertNoErr(err)
	if len(events) != 1 {
		t.Fatalf("Scanning Assets contract with 1 team returned %v events", len(events))
	}
}
