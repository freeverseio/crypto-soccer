package match_test

import (
	"math/big"
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

func TestMatchPlayer1stHalfWithDummyPlayers(t *testing.T) {
	bc, err := testutils.NewBlockchainNode()
	if err != nil {
		t.Fatal(err)
	}
	bc.DeployContracts(bc.Owner)

	m, err := match.NewMatch(bc.Contracts)
	if err != nil {
		t.Fatal(err)
	}
	m.MatchSeed = big.NewInt(54534)
	m.StartTime = big.NewInt(4343)
	for i := range m.HomeTeam.Players {
		m.HomeTeam.Players[i] = match.NewPlayer("3618502788719198628960202363453204454907735104658619445583958966799643443200")
	}
	for i := range m.VisitorTeam.Players {
		m.VisitorTeam.Players[i] = match.NewPlayer("3618502788732362665418772058558496602957087291216395979430084234696273690624")
	}

	err = m.Play1stHalf(bc.Contracts)
	if err != nil {
		t.Fatal(err)
	}

	if m.HomeGoals != 0 || m.VisitorGoals != 0 {
		t.Fatalf("Wrong result %v - %v", m.HomeGoals, m.VisitorGoals)
	}
}
