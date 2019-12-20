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
	team.State.Owner = "ciao"
	team.State.LeagueIdx = leagueIdx
	timezone := storage.Timezone{timezoneIdx}
	timezone.TimezoneCreate(tx)
	country := storage.Country{timezone.TimezoneIdx, countryIdx}
	country.CountryCreate(tx)
	league := storage.League{timezoneIdx, countryIdx, leagueIdx}
	league.LeagueCreate(tx)
	team.TeamCreate(tx)

	var player storage.Player
	player.PlayerId = big.NewInt(33)
	player.State.TeamId = big.NewInt(10)
	err = player.PlayerCreate(tx)
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
	team.State.Owner = "ciao"
	team.State.LeagueIdx = leagueIdx
	timezone := storage.Timezone{timezoneIdx}
	timezone.TimezoneCreate(tx)
	country := storage.Country{timezone.TimezoneIdx, countryIdx}
	country.CountryCreate(tx)
	league := storage.League{timezoneIdx, countryIdx, leagueIdx}
	league.LeagueCreate(tx)
	team.TeamCreate(tx)

	var player storage.Player
	player.PlayerId = big.NewInt(33)
	player.State.TeamId = big.NewInt(10)
	player.State.Name = "Iam Awesome"
	player.State.EncodedSkills = big.NewInt(4)
	err = player.PlayerCreate(tx)
	if err != nil {
		t.Fatal(err)
	}
	player2, err := storage.GetPlayer(tx, player.PlayerId)
	if err != nil {
		t.Fatal(err)
	}
	if player2.State.EncodedSkills.String() != player.State.EncodedSkills.String() {
		t.Fatal("Skills are different")
	}
	player2.State.EncodedSkills = big.NewInt(3)
	player2.State.RedCardMatchesLeft = 1
	player2.State.InjuryMatchesLeft = 3
	player2.State.Name = "Iam Sad"
	err = player2.PlayerUpdate(tx, player2.PlayerId, player2.State)
	if err != nil {
		t.Fatal(err)
	}
	player3, err := storage.GetPlayer(tx, player.PlayerId)
	if err != nil {
		t.Fatal(err)
	}
	if player2.State.RedCardMatchesLeft != player3.State.RedCardMatchesLeft {
		t.Fatal("Wrong RedCard")
	}
	if player2.State.InjuryMatchesLeft != player3.State.InjuryMatchesLeft {
		t.Fatal("Wrong InjuryMatchesLeft")
	}
	if player2.State.EncodedSkills.String() != player3.State.EncodedSkills.String() {
		t.Fatal("Skills player 3 are different")
	}
	if player3.State.Name != "Iam Sad" {
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
	team.State.Owner = "ciao"
	team.State.LeagueIdx = leagueIdx
	timezone := storage.Timezone{timezoneIdx}
	timezone.TimezoneCreate(tx)
	country := storage.Country{timezone.TimezoneIdx, countryIdx}
	country.CountryCreate(tx)
	league := storage.League{timezoneIdx, countryIdx, leagueIdx}
	league.LeagueCreate(tx)
	team.TeamCreate(tx)

	team.TeamID = big.NewInt(11)
	team.TeamCreate(tx)
	var player storage.Player
	player.PlayerId = big.NewInt(1)
	player.State.Defence = 4
	player.State.Endurance = 5
	player.State.Pass = 6
	player.State.Shoot = 7
	player.State.Speed = 8
	player.State.TeamId = big.NewInt(10)
	player.State.EncodedSkills, _ = new(big.Int).SetString("3618502788692870556043062973242620158809030731543066377891708431006382948352", 10)
	player.State.EncodedState, _ = new(big.Int).SetString("614878739568587161270510773682668741239185861458610514677961004951428661248", 10)

	err = player.PlayerCreate(tx)
	if err != nil {
		t.Fatal(err)
	}
	result, err := storage.GetPlayer(tx, player.PlayerId)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Equal(player) {
		t.Fatalf("Expected %v got %v", player, result)
	}
	player.State.Defence = 6
	player.State.TeamId = big.NewInt(11)
	err = player.PlayerUpdate(tx, player.PlayerId, player.State)
	if err != nil {
		t.Fatal(err)
	}
	result, err = storage.GetPlayer(tx, player.PlayerId)
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

	players, err := storage.GetPlayersOfTeam(tx, big.NewInt(343))
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
	team.State.Owner = "ciao"
	team.State.LeagueIdx = leagueIdx
	timezone := storage.Timezone{timezoneIdx}
	timezone.TimezoneCreate(tx)
	country := storage.Country{timezone.TimezoneIdx, countryIdx}
	country.CountryCreate(tx)
	league := storage.League{timezoneIdx, countryIdx, leagueIdx}
	league.LeagueCreate(tx)
	team.TeamCreate(tx)
	var player storage.Player
	player.PlayerId = big.NewInt(1)
	player.State.Defence = 4
	player.State.Endurance = 5
	player.State.Pass = 6
	player.State.Shoot = 7
	player.State.Speed = 8
	player.State.TeamId = team.TeamID
	player.State.EncodedSkills = big.NewInt(43535453)
	player.State.EncodedState = big.NewInt(43453)
	err = player.PlayerCreate(tx)
	if err != nil {
		t.Fatal(err)
	}
	player2 := player
	player2.PlayerId = big.NewInt(2)
	player2.State.EncodedSkills = big.NewInt(767)
	err = player2.PlayerCreate(tx)
	if err != nil {
		t.Fatal(err)
	}
	players, err = storage.GetPlayersOfTeam(tx, team.TeamID)
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
