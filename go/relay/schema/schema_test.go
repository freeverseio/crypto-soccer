package schema_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/relay/schema"
	graphql "github.com/graph-gophers/graphql-go"
	"gotest.tools/assert"
)

func TestTransferFirstBot(t *testing.T) {
	_, err := graphql.ParseSchema(schema.Schema, &schema.Resolver{})
	assert.NilError(t, err)
}
