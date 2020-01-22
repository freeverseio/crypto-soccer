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
	timezone.Insert(tx)
	country.Insert(tx)
	league.Insert(tx)
	team := storage.Team{}
	team.TeamID = big.NewInt(3)
	team.TimezoneIdx = timezone.TimezoneIdx
	team.CountryIdx = countryIdx
	team.Owner = "ciao"
	team.LeagueIdx = leagueIdx
	team.RankingPoints = math.MaxUint64
	if err = team.Insert(tx); err != nil {
		t.Fatal(err)
	}
	result, err := storage.TeamByTeamId(tx, team.TeamID)
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
	timezone.Insert(tx)
	country.Insert(tx)
	league.Insert(tx)

	var team storage.Team
	team.TeamID = big.NewInt(4)
	team.TimezoneIdx = timezone.TimezoneIdx
	team.CountryIdx = countryIdx
	team.Owner = "ciao"
	team.LeagueIdx = leagueIdx
	err = team.Insert(tx)
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
	teamResult, err := storage.TeamByTeamId(tx, team.TeamID)
	if err != nil {
		t.Fatal(err)
	}
	if teamResult.PrevPerfPoints != 0 {
		t.Fatalf("Wrong ranking points %v", teamResult.PrevPerfPoints)
	}
	if teamResult.RankingPoints != 0 {
		t.Fatalf("Wrong ranking points %v", teamResult.RankingPoints)
	}
}

func TestGetTeamOfUnexistenTeamID(t *testing.T) {
	tx, err := s.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	teamID := big.NewInt(434)
	_, err = storage.TeamByTeamId(tx, teamID)
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
	timezone.Insert(tx)
	country.Insert(tx)
	league.Insert(tx)

	var team storage.Team
	team.TeamID = big.NewInt(11)
	team.TimezoneIdx = timezone.TimezoneIdx
	team.CountryIdx = countryIdx
	team.Owner = "ciao"
	team.LeagueIdx = leagueIdx
	team.Insert(tx)
	teams, err := storage.TeamsByTimezoneIdxCountryIdxLeagueIdx(tx, timezone.TimezoneIdx, countryIdx, leagueIdx)
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
	timezone.Insert(tx)
	country.Insert(tx)
	league.Insert(tx)

	var team storage.Team
	team.TeamID = big.NewInt(4)
	team.TimezoneIdx = timezone.TimezoneIdx
	team.CountryIdx = countryIdx
	team.Owner = "ciao"
	team.LeagueIdx = leagueIdx
	team.RankingPoints = math.MaxUint64
	err = team.Insert(tx)
	if err != nil {
		t.Fatal(err)
	}
	team.Owner = "pippo"
	team.TrainingPoints = 4
	err = team.Update(tx)
	if err != nil {
		t.Fatal(err)
	}
	result, err := storage.TeamByTeamId(tx, team.TeamID)
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
// 	_, err = sto.TeamByTeamId(1)
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
// 	result, err := sto.TeamByTeamId(team.Id)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if result != team {
// 		t.Fatalf("Expected %v got %v", team, result)
// 	}
// }
