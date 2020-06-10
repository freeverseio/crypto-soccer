package process_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/synchronizer/process"
	"gotest.tools/assert"
)

func TestScanningIniting(t *testing.T) {
	scanner := process.NewEventScanner(bc.Contracts)
	assert.NilError(t, scanner.Process(nil))
	events := scanner.Events
	assert.Equal(t, 5, len(events))
	assert.Equal(t, events[0].Name, "ProxyNewDirectory")
	assert.Equal(t, events[1].Name, "AssetsInit")
	assert.Equal(t, events[2].Name, "AssetsTeamTransfer")
	assert.Equal(t, events[3].Name, "AssetsDivisionCreation")
	assert.Equal(t, events[4].Name, "AssetsTeamTransfer")
}

// func TestScanningTeamTransfer(t *testing.T) {
// 	eventCount := 0

// 	if scanner := process.NewEventScanner(bc.Contracts); scanner != nil {
// 		err := scanner.Process(nil)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		eventCount += len(scanner.Events)
// 		if eventCount == 0 {
// 			// skipping check for exact number at creation
// 			t.Fatalf("Expected some events received %v", eventCount)
// 		}
// 	} else {
// 		t.Fatal("Unable to create scanner")
// 	}

// 	timezoneIdx := uint8(1)
// 	countryIdx := big.NewInt(0)
// 	address := crypto.PubkeyToAddress(bc.Owner.PublicKey)
// 	tx, err := bc.Contracts.Assets.TransferFirstBotToAddr(bind.NewKeyedTransactor(bc.Owner), timezoneIdx, countryIdx, address)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	timezoneIdx = uint8(2)
// 	tx1, err := bc.Contracts.Assets.TransferFirstBotToAddr(bind.NewKeyedTransactor(bc.Owner), timezoneIdx, countryIdx, address)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	helper.WaitReceipt(bc.Client, tx, 3)
// 	helper.WaitReceipt(bc.Client, tx1, 3)

// 	if scanner := process.NewEventScanner(bc.Contracts); scanner != nil {
// 		err = scanner.Process(nil)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		if len(scanner.Events) != eventCount+2 {
// 			t.Fatalf("Expected 2 received %v", len(scanner.Events)-eventCount)
// 		}
// 	} else {
// 		t.Fatal("Unable to create scanner")
// 	}
// }
