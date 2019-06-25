package assets

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/testutils"
)

func TestScanTeamCreatedEmplyContract(t *testing.T) {
	blockchain, auth := testutils.InitBlockchain()

	statesContractAddress := common.HexToAddress("0x83a909262608c650bd9b0ae06e29d90d0f67ac5e")
	//Deploy contract
	_, _, contract, err := DeployAssets(
		auth,
		blockchain,
		statesContractAddress,
	)
	if err != nil {
		t.Fatal(err)
	}
	blockchain.Commit()

	events, err := contract.ScanTeamCreated()
	if err != nil {
		t.Fatal("Scanning error: ", err)
	}
	if len(events) != 0 {
		t.Fatalf("Scanning empty Assets contract returned %v events", len(events))
	}
}

func TestScanTeamCreated1TeamCreated(t *testing.T) {
	blockchain, auth := testutils.InitBlockchain()

	statesContractAddress := common.HexToAddress("0x83a909262608c650bd9b0ae06e29d90d0f67ac5e")
	//Deploy contract
	_, _, contract, err := DeployAssets(
		auth,
		blockchain,
		statesContractAddress,
	)
	if err != nil {
		t.Fatal(err)
	}
	blockchain.Commit()

	_, err = contract.CreateTeam(auth, "Barca", common.HexToAddress("0x83a909262608c650bd9b0ae06e29d90d0f67ac5e"))
	if err != nil {
		t.Fatal("Error creating team: ", err)
	}
	blockchain.Commit()

	events, err := contract.ScanTeamCreated()
	if err != nil {
		t.Fatal("Scanning error: ", err)
	}
	if len(events) != 1 {
		t.Fatalf("Scanning Assets contract with 1 team returned %v events", len(events))
	}
}
