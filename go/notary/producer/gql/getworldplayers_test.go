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

	assert.Equal(t, string(players[0].PlayerId()), "57896044618658097711785540140530465520588209856388886192059147693882236014541")
}
