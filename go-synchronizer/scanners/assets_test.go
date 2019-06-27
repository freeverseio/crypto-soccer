package scanners

import (
	"testing"

	//	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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

/*
// Before executing this test run scanners/run_ganache.sh
func TestGanache(t *testing.T) {
	ganache := testutils.NewGanache()
	bob := ganache.CreateAccountWithBalance("1000000000000000000") // 1 eth
	balance := ganache.GetBalance(testutils.CommonAddressFromPrivateKey(bob))
	if balance.String() != "1000000000000000000" {
		t.Fatalf("Funds were not transferred correctly. Transferred balance is %v", balance.Int64())
	}
	blocknum := ganache.GetLastBlockNumber()
	if blocknum == 0 {
		t.Fatalf("Block number should not be zero")
	}
	ganache.DeployContracts(bob)

	alice := ganache.CreateAccountWithBalance("50000000000000000000") // 50 eth
	ganache.CreateTeam("Barca", alice)

	events, err := ScanTeamCreated(ganache.Assets, nil)
	testutils.AssertNoErr(err)
	if len(events) != 1 {
		t.Fatalf("Scanning Assets contract with 1 team returned %v events", len(events))
	}

	events2, err2 := ScanTeamCreated(ganache.Assets, &bind.FilterOpts{Start: uint64(ganache.GetLastBlockNumber() + 1)})
	testutils.AssertNoErr(err2)
	if len(events2) != 0 {
		t.Fatalf("No new events should have been received, but got %v", len(events2))
	}

}
*/
