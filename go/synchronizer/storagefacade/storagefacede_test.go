package storagefacade_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/synchronizer/storagefacade"
	"gotest.tools/assert"
)

func TestNewMatchByLeagueWithNoMatches(t *testing.T) {
	tx, err := s.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	timezoneIdx := uint8(1)
	countryIdx := uint32(0)
	leagueIdx := uint32(0)
	day := uint8(0)
	matches, err := storagefacade.NewMatchesByLeague(tx, timezoneIdx, countryIdx, leagueIdx, day)
	assert.NilError(t, err)
	assert.Equal(t, len(matches), 0)
}

func TestNewMatchByLeagueWithMatches(t *testing.T) {
	tx, err := s.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	storagefacade.CreateMatch(t, tx)
	timezoneIdx := uint8(1)
	countryIdx := uint32(0)
	leagueIdx := uint32(0)
	day := uint8(0)
	matches, err := storagefacade.NewMatchesByLeague(tx, timezoneIdx, countryIdx, leagueIdx, day)
	assert.NilError(t, err)
	assert.Equal(t, len(matches), 1)
}
