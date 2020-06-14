package gql_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/relay/producer/gql"
	"github.com/freeverseio/crypto-soccer/go/relay/producer/gql/input"
	"gotest.tools/assert"
)

func TestSetTeamManagerNameNoTeamId(t *testing.T) {
	t.Parallel()
	in := input.SetTeamManagerNameInput{}
	resolver := gql.NewResolver(make(chan interface{}, 10))
	result, err := resolver.SetTeamManagerName(struct{ Input input.SetTeamManagerNameInput }{in})
	assert.NilError(t, err)
	assert.Equal(t, string(result), "")
}

func TestSetTeamManagerNameWithTeamId(t *testing.T) {
	t.Parallel()
	in := input.SetTeamManagerNameInput{}
	in.TeamId = "43534"
	resolver := gql.NewResolver(make(chan interface{}, 10))
	result, err := resolver.SetTeamManagerName(struct{ Input input.SetTeamManagerNameInput }{in})
	assert.NilError(t, err)
	assert.Equal(t, result, in.TeamId)
}
