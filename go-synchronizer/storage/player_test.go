package storage_test

import (
	"strings"
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

func TestPlayerAdd(t *testing.T) {
	sto, err := storage.NewSqlite3("../sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	var player storage.Player
	player.Id = 3
	player.State = "43524"
	err = sto.PlayerAdd(&player)
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
	if player != nil {
		t.Fatal("expected nil player")
	}
	id := uint64(3)
	state := "43524"
	var player2 storage.Player
	player2.Id = 3
	player2.State = "43524"
	err = sto.PlayerAdd(&player2)
	if err != nil {
		t.Fatal(err)
	}
	player, err = sto.GetPlayer(3)
	if err != nil {
		t.Fatal(err)
	}
	if player == nil {
		t.Fatal("expected player")
	}
	if player.Id != id {
		t.Fatalf("expected %v got %v", id, player.Id)
	}
	if strings.Compare(state, player.State) != 0 {
		t.Fatalf("Expected %v got %v", state, player.State)
	}
}
