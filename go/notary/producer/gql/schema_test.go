package gql_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql"
	graphql "github.com/graph-gophers/graphql-go"
	"gotest.tools/assert"
)

func TestSchemaParsing(t *testing.T) {
	_, err := graphql.ParseSchema(gql.Schema, nil)
	assert.NilError(t, err)
}
