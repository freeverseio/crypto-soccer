package storage_test

import (
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage"
)

// func TestTeamStateUpdate(t *testing.T) {
// 	sto, err := storage.NewSqlite3("../sql/00_schema.sql")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	team := storage.Team{
// 		TeamID:      big.NewInt(4),
// 		TimezoneIdx: 5,
// 		CountryIdx:  1,
// 		State: storage.TeamState{
// 			Owner: "io",
// 		},
// 	}
// 	err = sto.TeamCreate(team)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	team.State = storage.TeamState{
// 		BlockNumber:          6,
// 		Owner:                "tu",
// 		CurrentLeagueId:      6,
// 		PosInCurrentLeagueId: 3,
// 		PrevLeagueId:         1,
// 		PosInPrevLeagueId:    0,
// 	}
// 	err = sto.TeamStateUpdate(team.Id, team.State)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	result, err := sto.GetTeam(team.Id)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if team != result {
// 		t.Fatalf("Expected %v got %v", team, result)
// 	}
// }

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

func TestTeamCreate(t *testing.T) {
	sto, err := storage.NewSqlite3("../sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	timezone := uint8(1)
	countryIdx := uint32(4)
	leagueIdx := uint32(0)
	sto.TimezoneCreate(storage.Timezone{timezone})
	sto.CountryCreate(storage.Country{timezone, countryIdx})
	sto.LeagueCreate(storage.League{timezone, countryIdx, leagueIdx})
	var team storage.Team
	team.TeamID = big.NewInt(4)
	team.TimezoneIdx = timezone
	team.CountryIdx = countryIdx
	team.State.Owner = "ciao"
	team.State.LeagueIdx = leagueIdx
	err = sto.TeamCreate(team)
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

func TestGetTeamOfUnexistenTeamID(t *testing.T) {
	sto, err := storage.NewSqlite3("../sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	teamID := big.NewInt(434)
	_, err = sto.GetTeam(teamID)
	if err == nil {
		t.Fatal("Not error on unsexistent team")
	}
}

func TestGetTeamInLeague(t *testing.T) {
	sto, err := storage.NewSqlite3("../sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	timezone := uint8(1)
	countryIdx := uint32(0)
	leagueIdx := uint32(0)
	var team storage.Team
	team.TeamID = big.NewInt(11)
	team.TimezoneIdx = timezone
	team.CountryIdx = countryIdx
	team.State.Owner = "ciao"
	team.State.LeagueIdx = leagueIdx
	sto.TimezoneCreate(storage.Timezone{timezone})
	sto.CountryCreate(storage.Country{timezone, countryIdx})
	sto.LeagueCreate(storage.League{timezone, countryIdx, leagueIdx})
	sto.TeamCreate(team)
	teams, err := sto.GetTeamsInLeague(timezone, countryIdx, leagueIdx)
	if err != nil {
		t.Fatal(err)
	}
	if len(teams) != 1 {
		t.Fatalf("Expected 1 received %v", len(teams))
	}
}

func TestUpdateTeamOwner(t *testing.T) {
	sto, err := storage.NewSqlite3("../sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	timezone := uint8(1)
	countryIdx := uint32(4)
	leagueIdx := uint32(0)
	sto.TimezoneCreate(storage.Timezone{timezone})
	sto.CountryCreate(storage.Country{timezone, countryIdx})
	sto.LeagueCreate(storage.League{timezone, countryIdx, leagueIdx})
	var team storage.Team
	team.TeamID = big.NewInt(4)
	team.TimezoneIdx = timezone
	team.CountryIdx = countryIdx
	team.State.Owner = "ciao"
	team.State.LeagueIdx = leagueIdx
	err = sto.TeamCreate(team)
	if err != nil {
		t.Fatal(err)
	}
	team.State.Owner = "pippo"
	err = sto.TeamUpdate(team.TeamID, team.State)
	if err != nil {
		t.Fatal(err)
	}
	result, err := sto.GetTeam(team.TeamID)
	if err != nil {
		t.Fatal(err)
	}
	if team.State.Owner != result.State.Owner {
		t.Fatalf("expected owner pippo but got %v", result.State.Owner)
	}
}

// func TestTeamAddSameTimeTwice(t *testing.T) {
// 	sto, err := storage.NewSqlite3("../sql/00_schema.sql")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	var team storage.Team
// 	team.Id = 3
// 	team.Name = "ciao"
// 	err = sto.TeamAdd(team)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	err = sto.TeamAdd(team)
// 	if err == nil {
// 		t.Fatal("No error creating 2 teams with same id")
// 	}
// }

// func TestGetUnexistentTeam(t *testing.T) {
// 	sto, err := storage.NewSqlite3("../sql/00_schema.sql")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	_, err = sto.GetTeam(1)
// 	if err == nil {
// 		t.Fatal("Expecting error on unexistent team")
// 	}
// }

// func TestGetTeam(t *testing.T) {
// 	sto, err := storage.NewSqlite3("../sql/00_schema.sql")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	team := storage.Team{
// 		Id:                4,
// 		Name:              "pippo",
// 		CreationTimestamp: 67,
// 		CountryId:         1,
// 		State: storage.TeamState{
// 			BlockNumber:          5,
// 			Owner:                "io",
// 			CurrentLeagueId:      7,
// 			PosInCurrentLeagueId: 4,
// 			PrevLeagueId:         2,
// 			PosInPrevLeagueId:    1,
// 		},
// 	}

// 	err = sto.TeamAdd(team)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	result, err := sto.GetTeam(team.Id)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if result != team {
// 		t.Fatalf("Expected %v got %v", team, result)
// 	}
// }
