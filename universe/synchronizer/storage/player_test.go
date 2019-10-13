package storage_test

import (
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/universe/synchronizer/storage"

	_ "github.com/lib/pq"
)

func TestPlayerCount(t *testing.T) {
	storage, err := storage.NewSqlite3("../sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	count, err := storage.PlayerCount()
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("Expected 0 result %v", count)
	}
}

func TestPlayerCreate(t *testing.T) {
	sto, err := storage.NewSqlite3("../sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
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
	sto.TimezoneCreate(storage.Timezone{timezoneIdx})
	sto.CountryCreate(storage.Country{timezoneIdx, countryIdx})
	sto.LeagueCreate(storage.League{timezoneIdx, countryIdx, leagueIdx})
	sto.TeamCreate(team)
	var player storage.Player
	player.PlayerId = big.NewInt(33)
	player.State.TeamId = big.NewInt(10)
	err = sto.PlayerCreate(player)
	if err != nil {
		t.Fatal(err)
	}
	count, err := sto.PlayerCount()
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("Expected 1 result %v", count)
	}
}

func TestPlayerUpdate(t *testing.T) {
	sto, err := storage.NewSqlite3("../sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	var playerState storage.PlayerState
	err = sto.PlayerUpdate(big.NewInt(1), playerState)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetPlayer(t *testing.T) {
	sto, err := storage.NewSqlite3("../sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
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
	sto.TimezoneCreate(storage.Timezone{timezoneIdx})
	sto.CountryCreate(storage.Country{timezoneIdx, countryIdx})
	sto.LeagueCreate(storage.League{timezoneIdx, countryIdx, leagueIdx})
	sto.TeamCreate(team)
	team.TeamID = big.NewInt(11)
	sto.TeamCreate(team)
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

	err = sto.PlayerCreate(player)
	if err != nil {
		t.Fatal(err)
	}
	result, err := sto.GetPlayer(player.PlayerId)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Equal(player) {
		t.Fatalf("Expected %v got %v", player, result)
	}
	player.State.Defence = 6
	player.State.TeamId = big.NewInt(11)
	err = sto.PlayerUpdate(player.PlayerId, player.State)
	if err != nil {
		t.Fatal(err)
	}
	result, err = sto.GetPlayer(player.PlayerId)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Equal(player) {
		t.Fatalf("Expected %v got %v", player, result)
	}
}

func TestGetPlayersOfTeam(t *testing.T) {
	sto, err := storage.NewSqlite3("../sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	players, err := sto.GetPlayersOfTeam(big.NewInt(343))
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
	sto.TimezoneCreate(storage.Timezone{timezoneIdx})
	sto.CountryCreate(storage.Country{timezoneIdx, countryIdx})
	sto.LeagueCreate(storage.League{timezoneIdx, countryIdx, leagueIdx})
	sto.TeamCreate(team)
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
	err = sto.PlayerCreate(player)
	if err != nil {
		t.Fatal(err)
	}
	player2 := player
	player2.PlayerId = big.NewInt(2)
	player2.State.EncodedSkills = big.NewInt(767)
	err = sto.PlayerCreate(player2)
	if err != nil {
		t.Fatal(err)
	}
	players, err = sto.GetPlayersOfTeam(team.TeamID)
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

// func TestPlayerAddTwiceSameTeam(t *testing.T) {
// 	sto, err := storage.NewSqlite3("../sql/00_schema.sql")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	var player storage.Player
// 	player.Id = 3
// 	err = sto.PlayerAdd(player)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	err = sto.PlayerAdd(player)
// 	if err == nil {
// 		t.Fatal("No error adding the same player twice")
// 	}
// }
