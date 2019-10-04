package process_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go-synchronizer/process"
)

func TestScanningNothing(t *testing.T) {
	scanner := process.NewEventScanner()
	err := scanner.Process(nil)
	if err != nil {
		t.Fatal(err)
	}
	events := scanner.Events
	if len(events) != 0 {
		t.Fatalf("Wrong events number: %v", len(scanner.Events))
	}
}
