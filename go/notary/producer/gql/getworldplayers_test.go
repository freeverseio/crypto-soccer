package gql_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"gotest.tools/assert"
)

func TestGetWorldPlayers(t *testing.T) {
	ch := make(chan interface{}, 10)
	r := gql.NewResolver(ch, *bc.Contracts, namesdb, googleCredentials)

	in := input.GetWorldPlayersInput{}
	in.Signature = "a67621b4763db406f404c4a600ce0e79ee50147c209e85d2f146f0d760c0a1ac2a213a06f702995cee279af1f588b55c9fa462b2e6a9502d25cede77ec690ced1c"
	in.TeamId = "274877906944"

	players, err := r.GetWorldPlayers(struct{ Input input.GetWorldPlayersInput }{in})
	assert.NilError(t, err)
	assert.Equal(t, len(players), 30)
	assert.Equal(t, players[0].Race() == "", false)
	assert.Equal(t, players[0].CountryOfBirth() == "", false)
}

func TestCreateWorldPlayerBatch(t *testing.T) {
	value := int64(3000)
	maxPotential := uint8(9)
	now := int64(1554940800) // first second of a week
	teamId := "274877906944"

	players, err := gql.CreateWorldPlayerBatch(
		*bc.Contracts,
		namesdb,
		value,
		maxPotential,
		teamId,
		now-2,
	)
	assert.NilError(t, err)
	assert.Equal(t, len(players), 30)
	assert.Equal(t, string(players[0].PlayerId()), "25753057320981211674441424157481453821747928514686071527706372")
	assert.Equal(t, players[0].ValidUntil(), "1554940800")
	assert.Equal(t, players[0].Name(), "Bogdan Gimenez")
	assert.Equal(t, players[0].Race() == "", false)
	assert.Equal(t, players[0].CountryOfBirth() == "", false)
	assert.Equal(t, players[0].Speed(), int32(3402))

	players, err = gql.CreateWorldPlayerBatch(
		*bc.Contracts,
		namesdb,
		value,
		maxPotential,
		teamId,
		now-1,
	)
	assert.NilError(t, err)
	assert.Equal(t, len(players), 30)
	assert.Equal(t, string(players[0].PlayerId()), "25753057320981211674441424157481453821747928514686071527706372")
	assert.Equal(t, players[0].ValidUntil(), "1554940800")
	assert.Equal(t, players[0].Race() == "", false)
	assert.Equal(t, players[0].CountryOfBirth() == "", false)

	players, err = gql.CreateWorldPlayerBatch(
		*bc.Contracts,
		namesdb,
		value,
		maxPotential,
		teamId,
		now,
	)
	assert.NilError(t, err)
	assert.Equal(t, len(players), 30)
	assert.Equal(t, string(players[0].PlayerId()), "25753057320981211674441424157482721472348156744087568230911748")
	assert.Equal(t, players[0].ValidUntil(), "1555545600")

	players, err = gql.CreateWorldPlayerBatch(
		*bc.Contracts,
		namesdb,
		value,
		maxPotential,
		teamId,
		now+1,
	)
	assert.NilError(t, err)
	assert.Equal(t, len(players), 30)
	assert.Equal(t, string(players[0].PlayerId()), "25753057320981211674441424157482721472348156744087568230911748")
	assert.Equal(t, players[0].ValidUntil(), "1555545600")
}
