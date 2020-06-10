package gql_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/relay/producer/gql"
	"github.com/freeverseio/crypto-soccer/go/relay/producer/gql/input"
	"gotest.tools/assert"
)

func TestSetTeamNameNoTeamId(t *testing.T) {
	t.Parallel()
	in := input.SetTeamNameInput{}
	resolver := gql.NewResolver(make(chan interface{}, 10))
	result, err := resolver.SetTeamName(struct{ Input input.SetTeamNameInput }{in})
	assert.NilError(t, err)
	assert.Equal(t, string(result), "")
}

func TestSetTeamNameWithTeamId(t *testing.T) {
	t.Parallel()
	in := input.SetTeamNameInput{}
	in.TeamId = "43534"
	resolver := gql.NewResolver(make(chan interface{}, 10))
	result, err := resolver.SetTeamName(struct{ Input input.SetTeamNameInput }{in})
	assert.NilError(t, err)
	assert.Equal(t, result, in.TeamId)
}
