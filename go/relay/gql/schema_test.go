package gql_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/relay/gql"
	graphql "github.com/graph-gophers/graphql-go"
	"gotest.tools/assert"
)

func TestTransferFirstBot(t *testing.T) {
	_, err := graphql.ParseSchema(gql.Schema, gql.NewResolver())
	assert.NilError(t, err)
}
