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

/*
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
	fmt.Println("block:", blocknum)
	ganache.DeployContracts(bob)

	alice := ganache.CreateAccountWithBalance("40000000000000000000") // 400 eth
	ganache.CreateTeam("Barca", alice)

	//events, err := ScanTeamCreated(ganache.Assets, nil)
	//testutils.AssertNoErr(err)
	//if len(events) != 1 {
	//	t.Fatalf("Scanning Assets contract with 1 team returned %v events", len(events))
	//}
}
*/
