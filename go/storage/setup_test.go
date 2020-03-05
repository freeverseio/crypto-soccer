package storage_test

import (
	"database/sql"
	"log"
	"math"
	"os"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/storage"
	"github.com/freeverseio/crypto-soccer/go/testutils"
	"gotest.tools/assert"
)

var s *sql.DB
var bc *testutils.BlockchainNode

func TestMain(m *testing.M) {
	var err error
	s, err = storage.New("postgres://freeverse:freeverse@localhost:5432/cryptosoccer?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	bc, err = testutils.NewBlockchainNode()
	if err != nil {
		log.Fatal(err)
	}
	bc.DeployContracts(bc.Owner)

	os.Exit(m.Run())
}

const timezoneIdx = uint8(1)
const countryIdx = uint32(0)
const leagueIdx = uint32(0)
const teamID = "1"

func createMinimumUniverse(t *testing.T, tx *sql.Tx) {
	timezone := storage.Timezone{timezoneIdx}
	assert.NilError(t, timezone.Insert(tx))

	country := storage.Country{timezone.TimezoneIdx, countryIdx}
	assert.NilError(t, country.Insert(tx))

	league := storage.League{timezone.TimezoneIdx, countryIdx, leagueIdx}
	assert.NilError(t, league.Insert(tx))

	team := storage.NewTeam()
	team.TeamID = teamID
	team.TimezoneIdx = timezone.TimezoneIdx
	team.CountryIdx = countryIdx
	team.LeagueIdx = leagueIdx
	team.Owner = "my team"
	team.RankingPoints = math.MaxUint64
	assert.NilError(t, team.Insert(tx))
}
