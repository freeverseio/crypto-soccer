package assets

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestScanTeamCreatedEmplyContract(t *testing.T) {
	//Setup simulated block chain
	var gasLimit uint64 = 8000029
	key, _ := crypto.GenerateKey()
	auth := bind.NewKeyedTransactor(key)
	alloc := make(core.GenesisAlloc)
	alloc[auth.From] = core.GenesisAccount{Balance: big.NewInt(1000000000)}
	blockchain := backends.NewSimulatedBackend(alloc, gasLimit)

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

	events := scanTeamCreated(contract)
	if len(events) != 0 {
		t.Fatalf("Scanning empty Assets contract returned %v events", len(events))
	}
}

func TestScanTeamCreated1TeamCreated(t *testing.T) {
	//Setup simulated block chain
	var gasLimit uint64 = 8000029
	key, _ := crypto.GenerateKey()
	auth := bind.NewKeyedTransactor(key)
	alloc := make(core.GenesisAlloc)
	alloc[auth.From] = core.GenesisAccount{Balance: big.NewInt(1000000000)}
	blockchain := backends.NewSimulatedBackend(alloc, gasLimit)

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

	tr := bind.TransactOpts{
		From:   auth.From,
		Signer: auth.Signer,
		// GasLimit: big.NewInt(3141592),
	}

	_, err = contract.CreateTeam(&tr, "Barca", common.HexToAddress("0x83a909262608c650bd9b0ae06e29d90d0f67ac5e"))
	if err != nil {
		t.Fatal("Error creating team: ", err)
	}

	events := scanTeamCreated(contract)
	if len(events) != 1 {
		t.Fatalf("Scanning Assets contract with 1 team returned %v events", len(events))
	}
}
