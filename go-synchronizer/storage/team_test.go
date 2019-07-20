package storage_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage"
)

func TestTeamStateAdd(t *testing.T) {
	sto, err := storage.NewSqlite3("../sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	team := storage.Team{
		4,
		"pippo",
		"cavolfiore",
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
	err = sto.TeamStateAdd(team.Id, team.State)
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

func TestGetUnexistentTeamState(t *testing.T) {
	sto, err := storage.NewSqlite3("../sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	_, err = sto.GetTeamState(1)
	if err == nil {
		t.Fatal("No error on unexistent team")
	}
}

func TestGetTeamState(t *testing.T) {
	sto, err := storage.NewSqlite3("../sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	var teamState storage.TeamState
	teamState.BlockNumber = 33
	teamState.CurrentLeagueId = 3
	teamState.Owner = "44"
	teamState.PosInCurrentLeagueId = 4
	teamState.PosInPrevLeagueId = 7
	teamState.PrevLeagueId = 9
	err = sto.TeamStateAdd(1, teamState)
	if err != nil {
		t.Fatal(err)
	}
	result, err := sto.GetTeamState(1)
	if err != nil {
		t.Fatal(err)
	}
	if result != teamState {
		t.Fatalf("Expected %v got %v", teamState, result)
	}
	teamState.BlockNumber = 35
	teamState.CurrentLeagueId = 3
	teamState.Owner = "44"
	teamState.PosInCurrentLeagueId = 4
	teamState.PosInPrevLeagueId = 7
	teamState.PrevLeagueId = 9
	err = sto.TeamStateAdd(1, teamState)
	if err != nil {
		t.Fatal(err)
	}
	result, err = sto.GetTeamState(1)
	if err != nil {
		t.Fatal(err)
	}
	if result != teamState {
		t.Fatalf("Expected %v got %v", teamState, result)
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
		"cavolfiore",
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
