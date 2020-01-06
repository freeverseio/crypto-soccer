package match_test

import (
	"math/big"
	"testing"

	match "github.com/freeverseio/crypto-soccer/go/synchronizer/process/match"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"
)

func TestDefaultValues(t *testing.T) {
	seed := [32]byte{0x0}
	startTime := big.NewInt(0)

	sMatch := storage.Match{}
	homeTeam := storage.Team{}
	visitorTeam := storage.Team{}
	homeTeamPlayers := []*storage.Player{}
	visitorTeamPlayers := []*storage.Player{}
	homeTeamTactic := big.NewInt(0)
	visitorTeamTactic := big.NewInt(0)

	if mp := match.NewMatch(
		bc.Contracts,
		seed,
		startTime,
		&sMatch,
		&homeTeam,
		&visitorTeam,
		homeTeamPlayers,
		visitorTeamPlayers,
		homeTeamTactic,
		visitorTeamTactic,
	); mp == nil {
		t.Fatal("New instance is nil")
	}
}

func TestDefaultValuesPlayGame(t *testing.T) {
	seed := [32]byte{0x0}
	startTime := big.NewInt(0)

	homeTeam := storage.Team{}
	homeTeam.TeamID = big.NewInt(1)
	visitorTeam := storage.Team{}
	visitorTeam.TeamID = big.NewInt(1)

	sMatch := storage.Match{}
	sMatch.HomeTeamID = homeTeam.TeamID
	sMatch.VisitorTeamID = visitorTeam.TeamID
	sMatch.HomeMatchLog = big.NewInt(0)
	sMatch.VisitorMatchLog = big.NewInt(0)
	sMatch.HomeGoals = new(uint8)
	sMatch.VisitorGoals = new(uint8)

	homeTeamPlayers := []*storage.Player{}
	visitorTeamPlayers := []*storage.Player{}
	homeTeamTactic, err := match.GetEncodedTacticAtVerse(bc.Contracts, homeTeam.TeamID, 1)
	if err != nil {
		t.Fatal(err)
	}
	visitorTeamTactic, err := match.GetEncodedTacticAtVerse(bc.Contracts, visitorTeam.TeamID, 1)
	if err != nil {
		t.Fatal(err)
	}

	mp := match.NewMatch(
		bc.Contracts,
		seed,
		startTime,
		&sMatch,
		&homeTeam,
		&visitorTeam,
		homeTeamPlayers,
		visitorTeamPlayers,
		homeTeamTactic,
		visitorTeamTactic,
	)

	is2ndHalf := false
	if _, err := mp.Process(is2ndHalf); err != nil {
		t.Fatal(err)
	}
	if *mp.Match.HomeGoals != 0 {
		t.Fatalf("Wrong home goals %v", *mp.Match.HomeGoals)
	}
	if *mp.Match.VisitorGoals != 0 {
		t.Fatalf("Wrong visitor goals %v", *mp.Match.VisitorGoals)
	}
	if mp.Match.HomeMatchLog.String() != "1645504557321206042155578968558872826709262232930097591983538176" {
		t.Fatalf("Wrong home match log %v", mp.Match.HomeMatchLog.String())
	}
	if mp.Match.HomeMatchLog.String() != "1645504557321206042155578968558872826709262232930097591983538176" {
		t.Fatalf("Wrong visitor match log %v", mp.Match.VisitorMatchLog.String())
	}

	is2ndHalf = true
	if _, err := mp.Process(is2ndHalf); err != nil {
		t.Fatal(err)
	}
	if *mp.Match.HomeGoals != 0 {
		t.Fatalf("Wrong home goals %v", *mp.Match.HomeGoals)
	}
	if *mp.Match.VisitorGoals != 0 {
		t.Fatalf("Wrong visitor goals %v", *mp.Match.VisitorGoals)
	}
	if mp.Match.HomeMatchLog.String() != "1656419124875239866305548088509031409417165694142269319542924505165856768" {
		t.Fatalf("Wrong home match log %v", mp.Match.HomeMatchLog.String())
	}
	if mp.Match.HomeMatchLog.String() != "1656419124875239866305548088509031409417165694142269319542924505165856768" {
		t.Fatalf("Wrong visitor match log %v", mp.Match.VisitorMatchLog.String())
	}
}
