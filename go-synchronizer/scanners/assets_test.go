package scanners

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go-synchronizer/testutils"
)

func TestEmplyContract(t *testing.T) {
	cryptosoccer := testutils.CryptosoccerNew(t)
	events, err := ScanTeamCreated(cryptosoccer.AssetsContract)
	if err != nil {
		t.Fatal("Scanning error: ", err)
	}
	if len(events) != 0 {
		t.Fatalf("Scanning empty Assets contract returned %v events", len(events))
	}
}

func TestCreateTeam(t *testing.T) {
	cryptosoccer := testutils.CryptosoccerNew(t)
	_, err := cryptosoccer.AssetsContract.CreateTeam(cryptosoccer.Opts, "Barca", cryptosoccer.Opts.From)
	if err != nil {
		t.Fatal("Error creating team: ", err)
	}
	cryptosoccer.Commit()

	events, err := ScanTeamCreated(cryptosoccer.AssetsContract)
	if err != nil {
		t.Fatal("Scanning error: ", err)
	}
	if len(events) != 1 {
		t.Fatalf("Scanning Assets contract with 1 team returned %v events", len(events))
	}
}
