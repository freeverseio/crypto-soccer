package storage_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage"
)

func TestTeamStateUpdate(t *testing.T) {
	sto, err := storage.NewSqlite3("../sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	team := storage.Team{
		4,
		"pippo",
		54,
		storage.TeamState{
			BlockNumber:          5,
			Owner:                "io",
			CurrentLeagueId:      7,
			PosInCurrentLeagueId: 4,
			PrevLeagueId:         2,
			PosInPrevLeagueId:    1,
		},
	}
	err = sto.TeamAdd(team)
	if err != nil {
		t.Fatal(err)
	}
	team.State = storage.TeamState{
		BlockNumber:          6,
		Owner:                "tu",
		CurrentLeagueId:      6,
		PosInCurrentLeagueId: 3,
		PrevLeagueId:         1,
		PosInPrevLeagueId:    0,
	}
	err = sto.TeamStateUpdate(team.Id, team.State)
	if err != nil {
		t.Fatal(err)
	}
	result, err := sto.GetTeam(team.Id)
	if err != nil {
		t.Fatal(err)
	}
	if team != result {
		t.Fatalf("Expected %v got %v", team, result)
	}
}

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

func TestTeamAddSameTimeTwice(t *testing.T) {
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
	err = sto.TeamAdd(team)
	if err == nil {
		t.Fatal("No error creating 2 teams with same id")
	}
}

func TestGetUnexistentTeam(t *testing.T) {
	sto, err := storage.NewSqlite3("../sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	_, err = sto.GetTeam(1)
	if err == nil {
		t.Fatal("Expecting error on unexistent team")
	}
}

func TestGetTeam(t *testing.T) {
	sto, err := storage.NewSqlite3("../sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	team := storage.Team{
		4,
		"pippo",
		67,
		storage.TeamState{
			BlockNumber:          5,
			Owner:                "io",
			CurrentLeagueId:      7,
			PosInCurrentLeagueId: 4,
			PrevLeagueId:         2,
			PosInPrevLeagueId:    1,
		},
	}

	err = sto.TeamAdd(team)
	if err != nil {
		t.Fatal(err)
	}
	result, err := sto.GetTeam(team.Id)
	if err != nil {
		t.Fatal(err)
	}
	if result != team {
		t.Fatalf("Expected %v got %v", team, result)
	}
}
