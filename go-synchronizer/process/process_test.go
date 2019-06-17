package process

import (
	"testing"
	log "github.com/sirupsen/logrus"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage/memory"
	"math/big"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

// TODO : this is pattern to mocking contract !!!!!!


type AssetsMock struct {

}

func AssetsMockNew() *AssetsMock {
	return &AssetsMock{}
}

	func (m *AssetsMock) CountTeams(opts *bind.CallOpts) (*big.Int, error) {
		return big.NewInt(1), nil
	}

func assert(t *testing.T,err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func XTestInitBase(sto storage.Storage, t *testing.T) {
	assetsContract := AssetsMockNew()
	
	Process(assetsContract, sto)

	count, err := sto.TeamCount()
	if err != nil {
		t.Fatal(err)
	}

	countTeams, err := assetsContract.CountTeams(nil)
	if err != nil {
		log.Fatalf("Failed to retrieve token name: %v", err)
	}

	if count != countTeams.Uint64() {
		t.Fatal("team in bc", countTeams, " team in Database", count)
	}
}

func TestSyncTeam(t *testing.T) {
	mem := memory.New()
	XTestInitBase(mem,t)
}