package process_test

import (
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/synchronizer/process"
	"gotest.tools/assert"
)

func TestPlayerStateChangeGeneratePlayer(t *testing.T) {
	playerId, _ := new(big.Int).SetString("57896044618658097728977029352596290682772831803419867568648239153975217095645", 10)
	state, _ := new(big.Int).SetString("24519655528918455736691326674010135", 10)
	teamId := big.NewInt(1099511627776)
	player, err := process.GeneratePlayerByPlayerIdAndState(bc.Contracts, namesdb, 0, playerId, teamId, state)
	assert.NilError(t, err)
	assert.Equal(t, player.Name, "Brahim Herschler")
}

// func TestNewSpecialPlayer(t *testing.T) {
// 	t.Parallel()
// 	tx, err := universedb.Begin()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer tx.Rollback()

// 	timezoneIdx := uint8(0)
// 	countryIdx := uint32(0)
// 	leagueIdx := uint32(0)
// 	var team storage.Team
// 	team.TeamID = "1"
// 	team.TimezoneIdx = timezoneIdx
// 	team.CountryIdx = countryIdx
// 	team.Owner = "ciao"
// 	team.LeagueIdx = leagueIdx
// 	timezone := storage.Timezone{timezoneIdx}
// 	timezone.Insert(tx)
// 	country := storage.Country{timezone.TimezoneIdx, countryIdx}
// 	country.Insert(tx)
// 	league := storage.League{timezoneIdx, countryIdx, leagueIdx}
// 	league.Insert(tx)
// 	team.Insert(tx)

// 	playerId, _ := new(big.Int).SetString("57896044618658097728977029352596290682772831803419867568648239153975217095645", 10)
// 	state, _ := new(big.Int).SetString("24519655528918455736691326674010135", 10)
// 	player, err := process.GeneratePlayerByPlayerIdAndState(bc.Contracts, namesdb, 0, team.TeamID, playerId, state)
// 	assert.NilError(t, err)

// 	event := market.MarketPlayerStateChange{}
// 	event.PlayerId = playerId
// 	event.State = state
// 	if err = process.ConsumePlayerStateChange(tx, bc.Contracts, namesdb, event); err != nil {
// 		t.Fatal(err)
// 	}
// 	result, err := storage.PlayerByPlayerId(tx, playerId)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if !result.Equal(*player) {
// 		t.Fatalf("Expected %v got %v", player, result)
// 	}
// }
