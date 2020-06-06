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
	assert.Equal(t, len(players), 25)
	assert.Assert(t, players[0].Race() != "")
	assert.Assert(t, players[0].CountryOfBirth() != "")
}
