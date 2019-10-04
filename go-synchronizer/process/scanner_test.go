package process_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go-synchronizer/process"

	"github.com/freeverseio/crypto-soccer/go-synchronizer/testutils"
)

func TestScanningNothing(t *testing.T) {
	ganache := testutils.NewGanache()
	ganache.DeployContracts(ganache.Owner)

	leagues := ganache.Leagues
	updates := ganache.Updates

	scanner := process.NewEventScanner(leagues, updates)
	err := scanner.Process(nil)
	if err != nil {
		t.Fatal(err)
	}
	events := scanner.Events
	if len(events) != 1 {
		t.Fatalf("Wrong events number: %v", len(scanner.Events))
	}
	if events[0].Name != "LeaguesDivisionCreation" {
		t.Fatalf("Wrong event name %v", events[0].Name)
	}
}
