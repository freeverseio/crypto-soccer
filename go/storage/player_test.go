package storage_test

import (
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/storage"
	"gotest.tools/assert"

	_ "github.com/lib/pq"
)

func TestPlayerCount(t *testing.T) {
	t.Parallel()
	tx, err := s.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	count, err := storage.PlayerCount(tx)
	assert.NilError(t, err)
	assert.Equal(t, count, uint64(0))
}

func TestPlayerCreate(t *testing.T) {
	t.Parallel()
	tx, err := s.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	createMinimumUniverse(t, tx)

	var player storage.Player
	player.PlayerId = big.NewInt(33)
	player.TeamId = teamID

	blockNumber := uint64(3)
	assert.NilError(t, player.Insert(tx, blockNumber))
	count, err := storage.PlayerCount(tx)
	assert.NilError(t, err)
	assert.Equal(t, count, uint64(1))
}

func TestPlayerUpdate(t *testing.T) {
	t.Parallel()
	tx, err := s.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	createMinimumUniverse(t, tx)

	var player storage.Player
	player.PlayerId = big.NewInt(33)
	player.TeamId = teamID
	player.EncodedSkills = big.NewInt(4)
	player.YellowCard1stHalf = true
	blockNumber := uint64(3)
	assert.NilError(t, player.Insert(tx, blockNumber))
	player2, err := storage.PlayerByPlayerId(tx, player.PlayerId)
	assert.NilError(t, err)
	assert.Equal(t, player2.EncodedSkills.String(), player.EncodedSkills.String())
	player2.EncodedSkills = big.NewInt(3)
	player2.RedCard = true
	player2.YellowCard1stHalf = true
	player2.InjuryMatchesLeft = 3
	assert.NilError(t, player2.Update(tx, blockNumber+1))
	player3, err := storage.PlayerByPlayerId(tx, player.PlayerId)
	assert.NilError(t, err)
	assert.Equal(t, player2.RedCard, player3.RedCard)
	assert.Equal(t, player2.YellowCard1stHalf, player3.YellowCard1stHalf)
	assert.Equal(t, player2.InjuryMatchesLeft, player3.InjuryMatchesLeft)
	assert.Equal(t, player2.EncodedSkills.String(), player3.EncodedSkills.String())
}

func TestPlayerGetPlayer(t *testing.T) {
	tx, err := s.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	timezoneIdx := uint8(1)
	countryIdx := uint32(4)
	leagueIdx := uint32(0)
	var team storage.Team
	team.TeamID = "10"
	team.TimezoneIdx = timezoneIdx
	team.CountryIdx = countryIdx
	team.Owner = "ciao"
	team.LeagueIdx = leagueIdx
	timezone := storage.Timezone{timezoneIdx}
	timezone.Insert(tx)
	country := storage.Country{timezone.TimezoneIdx, countryIdx}
	country.Insert(tx)
	league := storage.League{timezoneIdx, countryIdx, leagueIdx}
	league.Insert(tx)
	team.Insert(tx)

	team.TeamID = "11"
	team.Insert(tx)
	var player storage.Player
	player.PlayerId = big.NewInt(1)
	player.Defence = 4
	player.Endurance = 5
	player.Pass = 6
	player.Shoot = 7
	player.Speed = 8
	player.TeamId = "10"
	player.EncodedSkills, _ = new(big.Int).SetString("3618502788692870556043062973242620158809030731543066377891708431006382948352", 10)
	player.EncodedState, _ = new(big.Int).SetString("614878739568587161270510773682668741239185861458610514677961004951428661248", 10)

	err = player.Insert(tx, uint64(3))
	if err != nil {
		t.Fatal(err)
	}
	result, err := storage.PlayerByPlayerId(tx, player.PlayerId)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Equal(player) {
		t.Fatalf("Expected %v got %v", player, result)
	}
	player.Defence = 6
	player.TeamId = "11"
	err = player.Update(tx, uint64(4))
	if err != nil {
		t.Fatal(err)
	}
	result, err = storage.PlayerByPlayerId(tx, player.PlayerId)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Equal(player) {
		t.Fatalf("Expected %v got %v", player, result)
	}
}

func TestPlayerGetPlayersOfTeam(t *testing.T) {
	tx, err := s.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	createMinimumUniverse(t, tx)

	players, err := storage.PlayersByTeamId(tx, "343")
	assert.NilError(t, err)
	assert.Equal(t, len(players), 0)
	var player storage.Player
	player.PlayerId = big.NewInt(1)
	player.Defence = 4
	player.Endurance = 5
	player.Pass = 6
	player.Shoot = 7
	player.Speed = 8
	player.TeamId = teamID
	player.EncodedSkills = big.NewInt(43535453)
	player.EncodedState = big.NewInt(43453)
	assert.NilError(t, player.Insert(tx, uint64(4)))
	player2 := player
	player2.PlayerId = big.NewInt(2)
	player2.EncodedSkills = big.NewInt(767)
	assert.NilError(t, player2.Insert(tx, uint64(4)))
	players, err = storage.PlayersByTeamId(tx, teamID)
	assert.NilError(t, err)
	assert.Equal(t, len(players), 2)
	assert.Assert(t, players[0].Equal(player))
	assert.Assert(t, players[1].Equal(player2))
}

func TestPlayerGetPlayersOfTeamDoNotCountDismissPlayers(t *testing.T) {
	tx, err := s.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	createMinimumUniverse(t, tx)

	players, err := storage.PlayersByTeamId(tx, "343")
	assert.NilError(t, err)
	assert.Equal(t, len(players), 0)
	var player storage.Player
	player.PlayerId = big.NewInt(1)
	player.Defence = 4
	player.Endurance = 5
	player.Pass = 6
	player.Shoot = 7
	player.Speed = 8
	player.TeamId = teamID
	player.EncodedSkills = big.NewInt(43535453)
	player.EncodedState = big.NewInt(43453)
	assert.NilError(t, player.Insert(tx, uint64(4)))
	player2 := player
	player2.PlayerId = big.NewInt(2)
	player2.EncodedSkills = big.NewInt(767)
	player2.ShirtNumber = 25
	assert.NilError(t, player2.Insert(tx, uint64(4)))
	players, err = storage.PlayersByTeamId(tx, teamID)
	assert.NilError(t, err)
	assert.Equal(t, len(players), 1)
	assert.Assert(t, players[0].Equal(player))
}
