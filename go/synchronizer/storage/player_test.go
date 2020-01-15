package storage_test

import (
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"

	_ "github.com/lib/pq"
)

func TestPlayerCount(t *testing.T) {
	tx, err := s.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	count, err := storage.PlayerCount(tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("Expected 0 result %v", count)
	}
}

func TestPlayerCreate(t *testing.T) {
	tx, err := s.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	timezoneIdx := uint8(1)
	countryIdx := uint32(4)
	leagueIdx := uint32(0)
	var team storage.Team
	team.TeamID = big.NewInt(10)
	team.TimezoneIdx = timezoneIdx
	team.CountryIdx = countryIdx
	team.Owner = "ciao"
	team.LeagueIdx = leagueIdx
	timezone := storage.Timezone{timezoneIdx}
	timezone.Insert(tx)
	country := storage.Country{timezone.TimezoneIdx, countryIdx}
	country.Insert(tx)
	league := storage.League{timezoneIdx, countryIdx, leagueIdx}
	league.Insert(tx)
	team.Insert(tx)

	var player storage.Player
	player.PlayerId = big.NewInt(33)
	player.TeamId = big.NewInt(10)
	err = player.Insert(tx)
	if err != nil {
		t.Fatal(err)
	}
	count, err := storage.PlayerCount(tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("Expected 1 result %v", count)
	}
}

func TestPlayerUpdate(t *testing.T) {
	tx, err := s.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	timezoneIdx := uint8(1)
	countryIdx := uint32(4)
	leagueIdx := uint32(0)
	var team storage.Team
	team.TeamID = big.NewInt(10)
	team.TimezoneIdx = timezoneIdx
	team.CountryIdx = countryIdx
	team.Owner = "ciao"
	team.LeagueIdx = leagueIdx
	timezone := storage.Timezone{timezoneIdx}
	timezone.Insert(tx)
	country := storage.Country{timezone.TimezoneIdx, countryIdx}
	country.Insert(tx)
	league := storage.League{timezoneIdx, countryIdx, leagueIdx}
	league.Insert(tx)
	team.Insert(tx)

	var player storage.Player
	player.PlayerId = big.NewInt(33)
	player.TeamId = big.NewInt(10)
	player.Name = "Iam Awesome"
	player.EncodedSkills = big.NewInt(4)
	err = player.Insert(tx)
	if err != nil {
		t.Fatal(err)
	}
	player2, err := storage.PlayerByPlayerId(tx, player.PlayerId)
	if err != nil {
		t.Fatal(err)
	}
	if player2.EncodedSkills.String() != player.EncodedSkills.String() {
		t.Fatal("Skills are different")
	}
	player2.EncodedSkills = big.NewInt(3)
	player2.RedCardMatchesLeft = 1
	player2.InjuryMatchesLeft = 3
	player2.Name = "Iam Sad"
	err = player2.Update(tx)
	if err != nil {
		t.Fatal(err)
	}
	player3, err := storage.PlayerByPlayerId(tx, player.PlayerId)
	if err != nil {
		t.Fatal(err)
	}
	if player2.RedCardMatchesLeft != player3.RedCardMatchesLeft {
		t.Fatal("Wrong RedCard")
	}
	if player2.InjuryMatchesLeft != player3.InjuryMatchesLeft {
		t.Fatal("Wrong InjuryMatchesLeft")
	}
	if player2.EncodedSkills.String() != player3.EncodedSkills.String() {
		t.Fatal("Skills player 3 are different")
	}
	if player3.Name != "Iam Sad" {
		t.Fatal("Wrong Name")
	}

}

func TestGetPlayer(t *testing.T) {
	tx, err := s.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	timezoneIdx := uint8(1)
	countryIdx := uint32(4)
	leagueIdx := uint32(0)
	var team storage.Team
	team.TeamID = big.NewInt(10)
	team.TimezoneIdx = timezoneIdx
	team.CountryIdx = countryIdx
	team.Owner = "ciao"
	team.LeagueIdx = leagueIdx
	timezone := storage.Timezone{timezoneIdx}
	timezone.Insert(tx)
	country := storage.Country{timezone.TimezoneIdx, countryIdx}
	country.Insert(tx)
	league := storage.League{timezoneIdx, countryIdx, leagueIdx}
	league.Insert(tx)
	team.Insert(tx)

	team.TeamID = big.NewInt(11)
	team.Insert(tx)
	var player storage.Player
	player.PlayerId = big.NewInt(1)
	player.Defence = 4
	player.Endurance = 5
	player.Pass = 6
	player.Shoot = 7
	player.Speed = 8
	player.TeamId = big.NewInt(10)
	player.EncodedSkills, _ = new(big.Int).SetString("3618502788692870556043062973242620158809030731543066377891708431006382948352", 10)
	player.EncodedState, _ = new(big.Int).SetString("614878739568587161270510773682668741239185861458610514677961004951428661248", 10)

	err = player.Insert(tx)
	if err != nil {
		t.Fatal(err)
	}
	result, err := storage.PlayerByPlayerId(tx, player.PlayerId)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Equal(player) {
		t.Fatalf("Expected %v got %v", player, result)
	}
	player.Defence = 6
	player.TeamId = big.NewInt(11)
	err = player.Update(tx)
	if err != nil {
		t.Fatal(err)
	}
	result, err = storage.PlayerByPlayerId(tx, player.PlayerId)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Equal(player) {
		t.Fatalf("Expected %v got %v", player, result)
	}
}

func TestGetPlayersOfTeam(t *testing.T) {
	tx, err := s.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	players, err := storage.PlayersByTeamId(tx, big.NewInt(343))
	if err != nil {
		t.Fatal(err)
	}
	if len(players) != 0 {
		t.Fatalf("Expected 0 received %v", len(players))
	}
	timezoneIdx := uint8(1)
	countryIdx := uint32(4)
	leagueIdx := uint32(0)
	var team storage.Team
	team.TeamID = big.NewInt(10)
	team.TimezoneIdx = timezoneIdx
	team.CountryIdx = countryIdx
	team.Owner = "ciao"
	team.LeagueIdx = leagueIdx
	timezone := storage.Timezone{timezoneIdx}
	timezone.Insert(tx)
	country := storage.Country{timezone.TimezoneIdx, countryIdx}
	country.Insert(tx)
	league := storage.League{timezoneIdx, countryIdx, leagueIdx}
	league.Insert(tx)
	team.Insert(tx)
	var player storage.Player
	player.PlayerId = big.NewInt(1)
	player.Defence = 4
	player.Endurance = 5
	player.Pass = 6
	player.Shoot = 7
	player.Speed = 8
	player.TeamId = team.TeamID
	player.EncodedSkills = big.NewInt(43535453)
	player.EncodedState = big.NewInt(43453)
	err = player.Insert(tx)
	if err != nil {
		t.Fatal(err)
	}
	player2 := player
	player2.PlayerId = big.NewInt(2)
	player2.EncodedSkills = big.NewInt(767)
	err = player2.Insert(tx)
	if err != nil {
		t.Fatal(err)
	}
	players, err = storage.PlayersByTeamId(tx, team.TeamID)
	if err != nil {
		t.Fatal(err)
	}
	if len(players) != 2 {
		t.Fatalf("Expected 2 received %v", len(players))
	}
	if !players[0].Equal(player) {
		t.Fatalf("Wrong player %v", players[0])
	}
	if !players[1].Equal(player2) {
		t.Fatalf("Wrong player %v", players[0])
	}
}
