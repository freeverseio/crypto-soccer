package match_test

import (
	"math/big"
	"testing"

	match "github.com/freeverseio/crypto-soccer/go/synchronizer/process/match"
)

func TestDefaultValues(t *testing.T) {
	seed := [32]byte{0x0}
	startTime := big.NewInt(0)

	homeTeam, _ := match.NewTeam(bc.Contracts)
	visitorTeam, _ := match.NewTeam(bc.Contracts)

	if mp := match.NewMatch(
		bc.Contracts,
		seed,
		startTime,
		homeTeam,
		visitorTeam,
	); mp == nil {
		t.Fatal("New instance is nil")
	}
}

func TestPlayi1stHalf(t *testing.T) {
	seed := [32]byte{0x0}
	startTime := big.NewInt(0)

	homeTeam, _ := match.NewTeam(bc.Contracts)
	visitorTeam, _ := match.NewTeam(bc.Contracts)
	homeTeam.TeamID = big.NewInt(1)
	visitorTeam.TeamID = big.NewInt(1)
	for i := 0; i < 25; i++ {
		homeTeam.Players[i] = match.NewPlayer("60912465658141224081372268432703414642709456376891023")
		visitorTeam.Players[i] = match.NewPlayer("60912465658141224081372268432703414642709456376891023")
	}

	mp := match.NewMatch(
		bc.Contracts,
		seed,
		startTime,
		homeTeam,
		visitorTeam,
	)

	is2ndHalf := false
	if _, err := mp.Process(is2ndHalf); err != nil {
		t.Fatal(err)
	}
	if mp.HomeGoals != 0 {
		t.Fatalf("Wrong home goals %v", mp.HomeGoals)
	}
	if mp.VisitorGoals != 0 {
		t.Fatalf("Wrong visitor goals %v", mp.VisitorGoals)
	}
	if mp.HomeMatchLog.String() != "754396374849259078542193939289572664399428455325249558202914101526528" {
		t.Fatalf("Wrong home match log %v", mp.HomeMatchLog.String())
	}
	if mp.VisitorMatchLog.String() != "754396374849259078542193939289572664399428455325249558202914101526528" {
		t.Fatalf("Wrong visitor match log %v", mp.VisitorMatchLog.String())
	}
}
