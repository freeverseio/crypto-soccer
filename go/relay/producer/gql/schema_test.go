package gql_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/relay/producer/gql"
	graphql "github.com/graph-gophers/graphql-go"
	"gotest.tools/assert"
)

func TestSchemaParsing(t *testing.T) {
	_, err := graphql.ParseSchema(gql.Schema, gql.NewResolver(nil, *bc.Contracts))
	assert.NilError(t, err)
}
