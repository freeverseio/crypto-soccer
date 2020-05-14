package gql_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql"
	"gotest.tools/assert"
)

func TestSubmitPlaystorePlayerPurchaseValidPlayer(t *testing.T) {
	value := int64(1000)     // TODO: value is forced to be 1000
	maxPotential := uint8(9) // TODO: value is forced to be 9
	teamId := "1099511627778"
	epoch := int64(1589456942)

	players, err := gql.CreateWorldPlayerBatch(
		*bc.Contracts,
		namesdb,
		value,
		maxPotential,
		teamId,
		epoch,
	)
	assert.NilError(t, err)

	ch := make(chan interface{}, 10)
	r := gql.NewResolver(ch, *bc.Contracts, namesdb, googleCredentials)
	isValid, err := r.IsValidPlayer(
		string(players[0].PlayerId()),
		value,
		maxPotential,
		teamId,
		epoch,
	)
	assert.NilError(t, err)
	assert.Assert(t, isValid)
}
