package testutils

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/assets"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/states"
)

type Cryptosoccer struct {
	Backend        *backends.SimulatedBackend
	Opts           *bind.TransactOpts
	StatesContract *states.States
	AssetsContract *assets.Assets
}

func CryptosoccerNew(t *testing.T) *Cryptosoccer {
	var gasLimit uint64 = 8000029
	key, _ := crypto.GenerateKey()
	opts := bind.NewKeyedTransactor(key)
	alloc := make(core.GenesisAlloc)
	alloc[opts.From] = core.GenesisAccount{Balance: big.NewInt(1000000000)}
	backend := backends.NewSimulatedBackend(alloc, gasLimit)

	//Deploy Assets contract
	statesContractAddress, _, statesContract, err := states.DeployStates(
		opts,
		backend,
	)
	if err != nil {
		t.Fatal(err)
	}

	//Deploy Assets contract
	_, _, assetsContract, err := assets.DeployAssets(
		opts,
		backend,
		statesContractAddress,
	)
	if err != nil {
		t.Fatal(err)
	}
	backend.Commit()

	return &Cryptosoccer{
		Backend:        backend,
		Opts:           opts,
		StatesContract: statesContract,
		AssetsContract: assetsContract,
	}
}
