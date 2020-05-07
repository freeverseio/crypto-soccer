package gql_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"gotest.tools/assert"
)

func TestSubmitPlayerPurchase(t *testing.T) {
	t.Skip("Rective this test to test google api")
	ch := make(chan interface{}, 10)
	r := gql.NewResolver(ch, *bc.Contracts, namesdb)

	in := input.SubmitPlayerPurchaseInput{}
	in.TeamId = "274877906944"
	in.PlayerId = "274877906944"
	in.PurchaseId = "korpimulxmslxissnschtkdb"
	in.Signature = "91366deb26195ac3b15b9e6fff99d425b2cc0d15d44dc8ee0377779400f92c4358a57754053facbe724e8a536e240b278cd651c756c46978eaebafc47767fd781b"

	id, err := r.SubmitPlayerPurchase(struct {
		Input input.SubmitPlayerPurchaseInput
	}{in})
	assert.NilError(t, err)
	assert.Equal(t, id, in.PlayerId)
}
