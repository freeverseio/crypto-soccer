package process_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/leagues"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/process"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/testutils"
)

func TestScanningNothing(t *testing.T) {
	scanner := process.NewEventScanner()
	err := scanner.Process()
	if err != nil {
		t.Fatal(err)
	}
	events := scanner.Events
	if len(events) != 0 {
		t.Fatalf("Wrong events number: %v", len(scanner.Events))
	}
}

func TestScanningIniting(t *testing.T) {
	ganache, err := testutils.NewBlockchainNode()
	if err != nil {
		t.Fatal(err)
	}
	err = ganache.DeployContracts(ganache.Owner)
	if err != nil {
		t.Fatal(err)
	}
	err = ganache.Init()
	if err != nil {
		t.Fatal(err)
	}

	divisionCreationIter, err := ganache.Leagues.FilterDivisionCreation(nil)
	if err != nil {
		t.Fatal(err)
	}
	scanner := process.NewEventScanner()
	err = scanner.ScanDivisionCreation(divisionCreationIter)
	if err != nil {
		t.Fatal(err)
	}
	err = scanner.Process()
	if err != nil {
		t.Fatal(err)
	}

	events := scanner.Events
	if len(events) != 24 {
		t.Fatalf("Expected 24 received %v", len(events))
	}
	for i := 0; i < 24; i++ {
		switch event := events[i].Value.(type) {
		case leagues.LeaguesDivisionCreation:
			expected := uint8(i + 1)
			if event.Timezone != expected {
				t.Fatalf("Expected %v recived %v", expected, event.Timezone)
			}
			break
		default:
			t.Fatalf("Wrong type %v", event)
		}
	}
}

func TestScanningTeamTransfer(t *testing.T) {
	ganache, err := testutils.NewBlockchainNode()
	if err != nil {
		t.Fatal(err)
	}
	ganache.DeployContracts(ganache.Owner)
	ganache.Init()

	scanner := process.NewEventScanner()

	iter, err := ganache.Leagues.FilterTeamTransfer(nil)
	if err != nil {
		t.Fatal(err)
	}
	err = scanner.ScanTeamTransfer(iter)
	if err != nil {
		t.Fatal(err)
	}
	err = scanner.Process()
	if err != nil {
		t.Fatal(err)
	}

	events := scanner.Events
	if len(events) != 0 {
		t.Fatalf("Expected 0 received %v", len(events))
	}

	timezoneIdx := uint8(1)
	countryIdx := big.NewInt(0)
	address := crypto.PubkeyToAddress(ganache.Owner.PublicKey)
	tx, err := ganache.Leagues.TransferFirstBotToAddr(bind.NewKeyedTransactor(ganache.Owner), timezoneIdx, countryIdx, address)
	if err != nil {
		t.Fatal(err)
	}
	timezoneIdx = uint8(2)
	tx1, err := ganache.Leagues.TransferFirstBotToAddr(bind.NewKeyedTransactor(ganache.Owner), timezoneIdx, countryIdx, address)
	if err != nil {
		t.Fatal(err)
	}
	ganache.WaitReceipt(tx, 3)
	ganache.WaitReceipt(tx1, 3)

	iter, err = ganache.Leagues.FilterTeamTransfer(nil)
	if err != nil {
		t.Fatal(err)
	}
	err = scanner.ScanTeamTransfer(iter)
	if err != nil {
		t.Fatal(err)
	}
	err = scanner.Process()
	if err != nil {
		t.Fatal(err)
	}
	events = scanner.Events
	if len(events) != 2 {
		t.Fatalf("Expected 2 received %v", len(events))
	}
}
