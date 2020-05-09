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
}

func TestCreateWorldPlayerBatch(t *testing.T) {
	value := int64(3000)
	now := int64(1554940800) // first second of a week
	teamId := "274877906944"

	players, err := gql.CreateWorldPlayerBatch(
		*bc.Contracts,
		namesdb,
		value,
		teamId,
		now-2,
	)
	assert.NilError(t, err)
	assert.Equal(t, len(players), 30)
	assert.Equal(t, string(players[0].PlayerId()), "57896044618658097711785532613796318893000562283106665962678893516101856333720")
	assert.Equal(t, players[0].ValidUntil(), "1554940800")
	assert.Equal(t, players[0].Name(), "Ekaitz Arana")
	assert.Equal(t, players[0].Speed(), int32(3413))

	players, err = gql.CreateWorldPlayerBatch(
		*bc.Contracts,
		namesdb,
		value,
		teamId,
		now-1,
	)
	assert.NilError(t, err)
	assert.Equal(t, len(players), 30)
	assert.Equal(t, string(players[0].PlayerId()), "57896044618658097711785532613796318893000562283106665962678893516101856333720")
	assert.Equal(t, players[0].ValidUntil(), "1554940800")

	players, err = gql.CreateWorldPlayerBatch(
		*bc.Contracts,
		namesdb,
		value,
		teamId,
		now,
	)
	assert.NilError(t, err)
	assert.Equal(t, len(players), 30)
	assert.Equal(t, string(players[0].PlayerId()), "57896044618658097711785532613796318893000562283106667171604713130731031039896")
	assert.Equal(t, players[0].ValidUntil(), "1555545600")

	players, err = gql.CreateWorldPlayerBatch(
		*bc.Contracts,
		namesdb,
		value,
		teamId,
		now+1,
	)
	assert.NilError(t, err)
	assert.Equal(t, len(players), 30)
	assert.Equal(t, string(players[0].PlayerId()), "57896044618658097711785532613796318893000562283106667171604713130731031039896")
	assert.Equal(t, players[0].ValidUntil(), "1555545600")
}
