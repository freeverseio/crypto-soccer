package storage_test

import (
	"math"
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"
)

func TestTeamCount(t *testing.T) {
	tx, err := s.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	count, err := storage.TeamCount(tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("Expected 0 result %v", count)
	}
}

func TestGetTeam(t *testing.T) {
	tx, err := s.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	countryIdx := uint32(4)
	leagueIdx := uint32(0)
	timezone := storage.Timezone{uint8(1)}
	country := storage.Country{timezone.TimezoneIdx, countryIdx}
	league := storage.League{timezone.TimezoneIdx, countryIdx, leagueIdx}
	timezone.TimezoneCreate(tx)
	country.CountryCreate(tx)
	league.LeagueCreate(tx)
	team := storage.Team{}
	team.TeamID = big.NewInt(3)
	team.TimezoneIdx = timezone.TimezoneIdx
	team.CountryIdx = countryIdx
	team.State.Owner = "ciao"
	team.State.LeagueIdx = leagueIdx
	team.State.RankingPoints = math.MaxUint64
	if err = team.TeamCreate(tx); err != nil {
		t.Fatal(err)
	}
	result, err := storage.GetTeam(tx, team.TeamID)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Equal(team) {
		t.Fatalf("Expected %v but %v", team, result)
	}
}

func TestTeamCreate(t *testing.T) {
	tx, err := s.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	countryIdx := uint32(4)
	leagueIdx := uint32(0)
	timezone := storage.Timezone{uint8(1)}
	country := storage.Country{timezone.TimezoneIdx, countryIdx}
	league := storage.League{timezone.TimezoneIdx, countryIdx, leagueIdx}
	timezone.TimezoneCreate(tx)
	country.CountryCreate(tx)
	league.LeagueCreate(tx)

	var team storage.Team
	team.TeamID = big.NewInt(4)
	team.TimezoneIdx = timezone.TimezoneIdx
	team.CountryIdx = countryIdx
	team.State.Owner = "ciao"
	team.State.LeagueIdx = leagueIdx
	err = team.TeamCreate(tx)
	if err != nil {
		t.Fatal(err)
	}
	count, err := storage.TeamCount(tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("Expected 1 result %v", count)
	}
	teamResult, err := storage.GetTeam(tx, team.TeamID)
	if err != nil {
		t.Fatal(err)
	}
	if teamResult.State.PrevPerfPoints != 0 {
		t.Fatalf("Wrong ranking points %v", teamResult.State.PrevPerfPoints)
	}
	if teamResult.State.RankingPoints != 0 {
		t.Fatalf("Wrong ranking points %v", teamResult.State.RankingPoints)
	}
}

func TestGetTeamOfUnexistenTeamID(t *testing.T) {
	tx, err := s.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	teamID := big.NewInt(434)
	_, err = storage.GetTeam(tx, teamID)
	if err == nil {
		t.Fatal("Not error on unsexistent team")
	}
}

func TestGetTeamInLeague(t *testing.T) {
	tx, err := s.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	countryIdx := uint32(1)
	leagueIdx := uint32(0)
	timezone := storage.Timezone{uint8(1)}
	country := storage.Country{timezone.TimezoneIdx, countryIdx}
	league := storage.League{timezone.TimezoneIdx, countryIdx, leagueIdx}
	timezone.TimezoneCreate(tx)
	country.CountryCreate(tx)
	league.LeagueCreate(tx)

	var team storage.Team
	team.TeamID = big.NewInt(11)
	team.TimezoneIdx = timezone.TimezoneIdx
	team.CountryIdx = countryIdx
	team.State.Owner = "ciao"
	team.State.LeagueIdx = leagueIdx
	team.TeamCreate(tx)
	teams, err := storage.GetTeamsInLeague(tx, timezone.TimezoneIdx, countryIdx, leagueIdx)
	if err != nil {
		t.Fatal(err)
	}
	if len(teams) != 1 {
		t.Fatalf("Expected 1 received %v", len(teams))
	}
}

func TestUpdateTeamOwner(t *testing.T) {
	tx, err := s.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	countryIdx := uint32(4)
	leagueIdx := uint32(0)
	timezone := storage.Timezone{uint8(1)}
	country := storage.Country{timezone.TimezoneIdx, countryIdx}
	league := storage.League{timezone.TimezoneIdx, countryIdx, leagueIdx}
	timezone.TimezoneCreate(tx)
	country.CountryCreate(tx)
	league.LeagueCreate(tx)

	var team storage.Team
	team.TeamID = big.NewInt(4)
	team.TimezoneIdx = timezone.TimezoneIdx
	team.CountryIdx = countryIdx
	team.State.Owner = "ciao"
	team.State.LeagueIdx = leagueIdx
	team.State.RankingPoints = math.MaxUint64
	err = team.TeamCreate(tx)
	if err != nil {
		t.Fatal(err)
	}
	team.State.Owner = "pippo"
	team.State.TrainingPoints = 4
	err = team.TeamUpdate(tx, team.TeamID, team.State)
	if err != nil {
		t.Fatal(err)
	}
	result, err := storage.GetTeam(tx, team.TeamID)
	if err != nil {
		t.Fatal(err)
	}
	if !team.Equal(result) {
		t.Fatalf("expected %v but got %v", team, result)
	}
}

// func TestTeamAddSameTimeTwice(t *testing.T) {
// 	sto, err := storage.NewSqlite3("../../../universe.db/00_schema.sql")
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
// 	sto, err := storage.NewSqlite3("../../../universe.db/00_schema.sql")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	_, err = sto.GetTeam(1)
// 	if err == nil {
// 		t.Fatal("Expecting error on unexistent team")
// 	}
// }

// func TestGetTeam(t *testing.T) {
// 	sto, err := storage.NewSqlite3("../../../universe.db/00_schema.sql")
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
