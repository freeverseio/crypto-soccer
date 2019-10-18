package process_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/contracts/leagues"
	"github.com/freeverseio/crypto-soccer/go/sync/process"
	"github.com/freeverseio/crypto-soccer/go/testutils"
)

func TestScanningNothing(t *testing.T) {
	scanner := process.NewEventScanner(nil, nil, nil)
	if scanner != nil {
		t.Fatal("scanner cannot be created with null contracts")
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

	scanner := process.NewEventScanner(ganache.Leagues, ganache.Updates, ganache.Market)
	err = scanner.Process(nil)
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

	eventCount := 0

	if scanner := process.NewEventScanner(ganache.Leagues, ganache.Updates, ganache.Market); scanner != nil {
		err = scanner.Process(nil)
		if err != nil {
			t.Fatal(err)
		}
		eventCount += len(scanner.Events)
		if eventCount == 0 {
			// skipping check for exact number at creation
			t.Fatalf("Expected some events received %v", eventCount)
		}
	} else {
		t.Fatal("Unable to create scanner")
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

	if scanner := process.NewEventScanner(ganache.Leagues, ganache.Updates, ganache.Market); scanner != nil {
		err = scanner.Process(nil)
		if err != nil {
			t.Fatal(err)
		}
		if len(scanner.Events) != eventCount+2 {
			t.Fatalf("Expected 2 received %v", len(scanner.Events)-eventCount)
		}
	} else {
		t.Fatal("Unable to create scanner")
	}
}
