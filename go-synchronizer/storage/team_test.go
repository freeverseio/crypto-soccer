package storage_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage"
)

func TestTeamCount(t *testing.T) {
	storage, err := storage.NewSqlite3("../sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	count, err := storage.TeamCount()
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("Expected 0 result %v", count)
	}
}

func TestTeamAdd(t *testing.T) {
	sto, err := storage.NewSqlite3("../sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	var team storage.Team
	team.Id = 3
	team.Name = "ciao"
	err = sto.TeamAdd(team)
	if err != nil {
		t.Fatal(err)
	}
	count, err := sto.TeamCount()
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("Expected 1 result %v", count)
	}
}

func TestGetTeam(t *testing.T) {
	sto, err := storage.NewSqlite3("../sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	_, err = sto.GetTeam(1)
	if err != nil {
		t.Fatal("Expecting nil")
	}
	var team storage.Team
	team.Id = 3
	team.Name = "ciao"
	err = sto.TeamAdd(team)
	if err != nil {
		t.Fatal(err)
	}
	result, err := sto.GetTeam(3)
	if err != nil {
		t.Fatal(err)
	}
	if result != team {
		t.Fatalf("Expected %v got %v", team, result)
	}
}
