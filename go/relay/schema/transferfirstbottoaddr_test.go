package schema_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/relay/schema"
	graphql "github.com/graph-gophers/graphql-go"
	"gotest.tools/assert"
)

func TestTransferFirstBotToAddr(t *testing.T) {
	schema := graphql.MustParseSchema(schema.Schema, schema.NewResolver())
	query := `mutation {
		transferFirstBotToAddr(timezone: 3, countryIdxInTimezone: "4", address: "0x0")
	}`
	ctx := context.Background()
	resp := schema.Exec(ctx, query, "", nil)
	json, err := json.Marshal(resp)
	assert.NilError(t, err)
	assert.Equal(t, string(json), `{"data":{"transferFirstBotToAddr":false}}`)
}
