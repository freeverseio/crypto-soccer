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
	err = sto.PlayerStateAdd(1, playerState)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetPlayerState(t *testing.T) {
	sto, err := storage.NewSqlite3("../sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	var playerState storage.PlayerState
	playerState.BlockNumber = "33"
	playerState.Defence = 4
	playerState.Endurance = 5
	playerState.Pass = 6
	playerState.Shoot = 7
	playerState.Speed = 8
	playerState.State = "23"
	playerState.TeamId = 99
	err = sto.PlayerStateAdd(1, playerState)
	if err != nil {
		t.Fatal(err)
	}
	result, err := sto.GetPlayerState(1)
	if err != nil {
		t.Fatal(err)
	}
	if result != playerState {
		t.Fatalf("Expected %v got %v", playerState, result)
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

func TestGetPlayer(t *testing.T) {
	sto, err := storage.NewSqlite3("../sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	player, err := sto.GetPlayer(1)
	if err != nil {
		t.Fatal(err)
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
	player, err = sto.GetPlayer(3)
	if err != nil {
		t.Fatal(err)
	}
	if player2 != player {
		t.Fatalf("Expected %v got %v", player2, player)
	}
}
