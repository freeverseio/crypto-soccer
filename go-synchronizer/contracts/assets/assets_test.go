package assets

import (
	"testing"
)

func TestDeployAssets(t *testing.T) {
	blockchain := DefaultSimulatedBlockchain()
	count := blockchain.CountTeams()
	if count.Int64() != 0 {
		t.Fatal("number of teams is not 0: ", count)
	}
}

func TestScanTeamCreatedEmptyContract(t *testing.T) {
	blockchain := DefaultSimulatedBlockchain()
	events := blockchain.ScanTeamCreated()
	if len(events) != 0 {
		t.Fatalf("Scanning empty Assets contract returned %v events", len(events))
	}
}

func TestScanTeamCreated1TeamCreated(t *testing.T) {
	blockchain := DefaultSimulatedBlockchain()
	alice := blockchain.CreateAccountWithBalance("1000000000000000000") // 1 eth
	blockchain.CreateTeam("Barca", alice)
	events := blockchain.ScanTeamCreated()
	if len(events) != 1 {
		t.Fatalf("Scanning Assets contract with 1 team returned %v events", len(events))
	}
}
