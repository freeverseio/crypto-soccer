package gql_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/relay/producer/gql"
	graphql "github.com/graph-gophers/graphql-go"
	"gotest.tools/assert"
)

func TestResolverTransferFirstBotToAddrResponse(t *testing.T) {
	schema := graphql.MustParseSchema(gql.Schema, gql.NewResolver(nil, *bc.Contracts))
	query := `mutation {
		transferFirstBotToAddr(timezone: 3, countryIdxInTimezone: "4", address: "0x0")
	}`
	ctx := context.Background()
	resp := schema.Exec(ctx, query, "", nil)
	json, err := json.Marshal(resp)
	assert.NilError(t, err)
	assert.Equal(t, string(json), `{"data":{"transferFirstBotToAddr":true}}`)
}
