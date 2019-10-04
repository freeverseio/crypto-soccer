package process_test

import (
	"testing"

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
	ganache := testutils.NewGanache()
	ganache.DeployContracts(ganache.Owner)
	ganache.Init()

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
}
