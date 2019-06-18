package assets

import (
	"log"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestDeplyAssets(t *testing.T) {
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

	count, err := contract.CountTeams(nil)
	if err != nil {
		log.Fatalf("Failed to count teams: %v", err)
	}
	if count.Int64() != 0 {
		t.Fatal("number of teams is not 0: ", count)
	}
}
