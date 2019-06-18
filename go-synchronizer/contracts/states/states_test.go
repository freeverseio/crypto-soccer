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

func TestSimulatedBackend(t *testing.T) {
	//Setup simulated block chain
	var gasLimit uint64 = 8000029
	key, _ := crypto.GenerateKey()
	auth := bind.NewKeyedTransactor(key)
	alloc := make(core.GenesisAlloc)
	alloc[auth.From] = core.GenesisAccount{Balance: big.NewInt(1000000000)}
	blockchain := backends.NewSimulatedBackend(alloc, gasLimit)

	//Deploy contract
	address, _, _, err := DeployStates(
		auth,
		blockchain,
	)
	if err != nil {
		t.Fatal(err)
	}

	// commit all pending transactions
	blockchain.Commit()

	log.Fatal(address)
}
