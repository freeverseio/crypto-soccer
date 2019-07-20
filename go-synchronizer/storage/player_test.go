package storage_test

import (
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

func TestPlayerStateAdd(t *testing.T) {
	sto, err := storage.NewSqlite3("../sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	var playerState storage.PlayerState
	err = sto.PlayerStateUpdate(1, playerState)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetPlayerState(t *testing.T) {
	sto, err := storage.NewSqlite3("../sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	var player storage.Player
	player.Id = 1
	player.MonthOfBirthInUnixTime = "ff"
	player.State.BlockNumber = 33
	player.State.Defence = 4
	player.State.Endurance = 5
	player.State.Pass = 6
	player.State.Shoot = 7
	player.State.Speed = 8
	player.State.State = "23"
	player.State.TeamId = 99
	err = sto.PlayerAdd(player)
	if err != nil {
		t.Fatal(err)
	}
	result, err := sto.GetPlayerState(player.Id)
	if err != nil {
		t.Fatal(err)
	}
	if result != player.State {
		t.Fatalf("Expected %v got %v", player.State, result)
	}
	player.State.BlockNumber = 366
	player.State.Defence = 6
	err = sto.PlayerStateUpdate(player.Id, player.State)
	if err != nil {
		t.Fatal(err)
	}
	result, err = sto.GetPlayerState(1)
	if err != nil {
		t.Fatal(err)
	}
	if result != player.State {
		t.Fatalf("Expected %v got %v", player.State, result)
	}
}

func TestPlayerAdd(t *testing.T) {
	sto, err := storage.NewSqlite3("../sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	var player storage.Player
	player.Id = 3
	err = sto.PlayerAdd(player)
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

func TestPlayerAddTwiceSameTeam(t *testing.T) {
	sto, err := storage.NewSqlite3("../sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	var player storage.Player
	player.Id = 3
	err = sto.PlayerAdd(player)
	if err != nil {
		t.Fatal(err)
	}
	err = sto.PlayerAdd(player)
	if err == nil {
		t.Fatal("No error adding the same player twice")
	}
}

func TestGetPlayer(t *testing.T) {
	sto, err := storage.NewSqlite3("../sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	player, err := sto.GetPlayer(1)
	if err == nil {
		t.Fatal("No error on get unexistent player")
	}
	var player2 storage.Player
	player2.Id = 3
	player2.MonthOfBirthInUnixTime = "4534"
	player2.State.Defence = 4
	// player2.State = "43524"
	err = sto.PlayerAdd(player2)
	if err != nil {
		t.Fatal(err)
	}
	player, err = sto.GetPlayer(player2.Id)
	if err != nil {
		t.Fatal(err)
	}
	if player2 != player {
		t.Fatalf("Expected %v got %v", player2, player)
	}
}
