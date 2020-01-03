package match_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/synchronizer/match"
	"github.com/freeverseio/crypto-soccer/go/testutils"
)

func TestMatchPlay1stHalf(t *testing.T) {
	bc, err := testutils.NewBlockchainNode()
	if err != nil {
		t.Fatal(err)
	}
	bc.DeployContracts(bc.Owner)

	m, err := match.NewMatch(bc.Contracts)
	if err != nil {
		t.Fatal(err)
	}

	err = m.Play1stHalf(bc.Contracts)
	if err != nil {
		t.Fatal(err)
	}

	if m.HomeGoals != 0 || m.VisitorGoals != 0 {
		t.Fatalf("Wrong result %v - %v", m.HomeGoals, m.VisitorGoals)
	}
}
