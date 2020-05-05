package gql_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"gotest.tools/assert"
)

func TestGetWorldPlayers(t *testing.T) {
	ch := make(chan interface{}, 10)
	r := gql.NewResolver(ch, *bc.Contracts)

	in := input.GetWorldPlayersInput{}
	in.Signature = "82b6568d3e792df067a07ca67316b916de3064ef0cdabcbf25a59e5e9745caa328ae510bd2a62a92e2f9710aa38798a0a7e7f47b0632bf08fa4c7abd52e5c0a11b"
	in.TeamId = "4"

	players, err := r.GetWorldPlayers(struct{ Input input.GetWorldPlayersInput }{in})
	assert.NilError(t, err)
	assert.Equal(t, len(players), 30)
}

func TestCreateWorldPlayerBatch(t *testing.T) {
	value := int64(3000)
	now := int64(3600 * 7)
	teamId := "6"

	players, err := gql.CreateWorldPlayerBatch(
		*bc.Contracts,
		value,
		teamId,
		now-2,
	)
	assert.NilError(t, err)
	assert.Equal(t, len(players), 30)
	assert.Equal(t, string(players[0].PlayerId()), "57896044618658097711785529983094086308552742237954859663381640645247386192918")
	assert.Equal(t, players[0].ValidUntil(), "20")

	players, err = gql.CreateWorldPlayerBatch(
		*bc.Contracts,
		value,
		teamId,
		now-1,
	)
	assert.NilError(t, err)
	assert.Equal(t, len(players), 30)
	assert.Equal(t, string(players[0].PlayerId()), "57896044618658097711785529983094086308552742237954859663381640645247386192918")

	players, err = gql.CreateWorldPlayerBatch(
		*bc.Contracts,
		value,
		teamId,
		now,
	)
	assert.NilError(t, err)
	assert.Equal(t, len(players), 30)
	assert.Equal(t, string(players[0].PlayerId()), "57896044618658097711785550502574930679863446151818835190092890767547663847502")

	players, err = gql.CreateWorldPlayerBatch(
		*bc.Contracts,
		value,
		teamId,
		now+1,
	)
	assert.NilError(t, err)
	assert.Equal(t, len(players), 30)
	assert.Equal(t, string(players[0].PlayerId()), "57896044618658097711785550502574930679863446151818835190092890767547663847502")
}
