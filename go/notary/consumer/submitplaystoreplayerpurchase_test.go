package consumer_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/consumer"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"gotest.tools/assert"
)

func TestSubmitPlayStorePlayerPurchaseEmptyEvent(t *testing.T) {
	in := input.SubmitPlayStorePlayerPurchaseInput{}
	assert.Error(t, consumer.SubmitPlayStorePlayerPurchase(
		*bc.Contracts,
		bc.Owner,
		googleCredentials,
		in,
		false,
	), "invalid playerId ")
}

func TestSubmitPlayStorePlayerPurchase(t *testing.T) {
	in := input.SubmitPlayStorePlayerPurchaseInput{}
	in.PlayerId = "3"
	in.TeamId = "4"
	in.PackageName = "PackageId"
	in.ProductId = "ProductId"
	in.PurchaseToken = "44434r"
	assert.Error(t, consumer.SubmitPlayStorePlayerPurchase(
		*bc.Contracts,
		bc.Owner,
		googleCredentials,
		in,
		false,
	), "unexpected end of JSON input")
}
