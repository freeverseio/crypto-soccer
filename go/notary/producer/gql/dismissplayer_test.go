package gql_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"gotest.tools/assert"
)

func TestDismissPlayer(t *testing.T) {
	ch := make(chan interface{}, 10)
	r := gql.NewResolver(ch, *bc.Contracts, namesdb, googleCredentials, service)

	in := input.DismissPlayerInput{}

	_, err := r.DismissPlayer(struct{ Input input.DismissPlayerInput }{in})
	assert.Error(t, err, "invalid validUntil")

	in.ValidUntil = "32323"
	_, err = r.DismissPlayer(struct{ Input input.DismissPlayerInput }{in})
	assert.Error(t, err, "invalid playerId")

	in.PlayerId = "32323"
	_, err = r.DismissPlayer(struct{ Input input.DismissPlayerInput }{in})
	assert.Error(t, err, "signature must be 65 bytes long")

	in.Signature = "0f13e4028d911bbf7e305267d593c6b67888030032e73f94a5cf8af204567ab629848e9290568aa5d19c1b7a4761a20ed4059072aacd79bde56e1b52c17a21311b"
	_, err = r.DismissPlayer(struct{ Input input.DismissPlayerInput }{in})
	assert.Error(t, err, "not player owner")
}
