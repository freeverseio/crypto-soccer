package storage_test

import (
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage"

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
	countryIdx := uint16(4)
	leagueIdx := uint8(0)
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

// func TestPlayerStateUpdate(t *testing.T) {
// 	sto, err := storage.NewSqlite3("../sql/00_schema.sql")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	var playerState storage.PlayerState
// 	err = sto.PlayerStateUpdate(1, playerState)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }

// func TestGetPlayer(t *testing.T) {
// 	sto, err := storage.NewSqlite3("../sql/00_schema.sql")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	var player storage.Player
// 	player.Id = 1
// 	player.MonthOfBirthInUnixTime = "ff"
// 	player.State.BlockNumber = 33
// 	player.State.Defence = 4
// 	player.State.Endurance = 5
// 	player.State.Pass = 6
// 	player.State.Shoot = 7
// 	player.State.Speed = 8
// 	player.State.State = "23"
// 	player.State.TeamId = 99
// 	err = sto.PlayerAdd(player)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	result, err := sto.GetPlayer(player.Id)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if result != player {
// 		t.Fatalf("Expected %v got %v", player, result)
// 	}
// 	player.State.BlockNumber = 366
// 	player.State.Defence = 6
// 	err = sto.PlayerStateUpdate(player.Id, player.State)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	result, err = sto.GetPlayer(player.Id)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if result != player {
// 		t.Fatalf("Expected %v got %v", player, result)
// 	}
// }

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
