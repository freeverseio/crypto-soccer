package states

import (
	"log"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestDeployStates(t *testing.T) {
	//Setup simulated block chain
	var gasLimit uint64 = 8000029
	key, _ := crypto.GenerateKey()
	auth := bind.NewKeyedTransactor(key)
	alloc := make(core.GenesisAlloc)
	alloc[auth.From] = core.GenesisAccount{Balance: big.NewInt(1000000000)}
	blockchain := backends.NewSimulatedBackend(alloc, gasLimit)

	//Deploy contract
	_, _, contract, err := DeployStates(
		auth,
		blockchain,
	)
	if err != nil {
		t.Fatal(err)
	}
	blockchain.Commit()

	isValid, err := contract.IsValidPlayerState(nil, big.NewInt(0))
	if err != nil {
		log.Fatal("Failed to call isValidPlayerState", err)
	}
	if isValid {
		t.Fatal("invalid state is valid")
	}
}
