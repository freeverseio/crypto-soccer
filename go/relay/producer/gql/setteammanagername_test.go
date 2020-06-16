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
	resolver := gql.NewResolver(make(chan interface{}, 10), *bc.Contracts)
	_, err := resolver.SetTeamManagerName(struct{ Input input.SetTeamManagerNameInput }{in})
	assert.Error(t, err, "Invalid TeamId")
}

func TestSetTeamManagerNameWithTeamId(t *testing.T) {
	t.Parallel()
	in := input.SetTeamManagerNameInput{}
	in.TeamId = "4"
	in.Name = "ciao"
	in.Signature = "3feac668bb718f492638b9b58d1f294379cdc8bde40074f5e49c3f80f28190e121f0fd08227c64a643dd032748ef772b0d1cf1500f649345521c133290c941a91b"
	resolver := gql.NewResolver(make(chan interface{}, 10), *bc.Contracts)
	_, err := resolver.SetTeamManagerName(struct{ Input input.SetTeamManagerNameInput }{in})
	assert.Error(t, err, "not allowed")
}
