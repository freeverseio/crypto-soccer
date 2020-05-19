package consumer_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/consumer"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"gotest.tools/assert"
)

func TestSubmitPlayStorePlayerPurchaseEmptyEvent(t *testing.T) {
	tx, err := db.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	in := input.SubmitPlayStorePlayerPurchaseInput{}
	assert.NilError(t, consumer.SubmitPlayStorePlayerPurchase(
		*bc.Contracts,
		tx,
		bc.Owner,
		googleCredentials,
		in,
		false,
	))
}

func TestSubmitPlayStorePlayerPurchase(t *testing.T) {
	tx, err := db.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	in := input.SubmitPlayStorePlayerPurchaseInput{}
	in.PlayerId = "3"
	in.TeamId = "4"
	in.Receipt = "PackageId"
	assert.NilError(t, consumer.SubmitPlayStorePlayerPurchase(
		*bc.Contracts,
		tx,
		bc.Owner,
		googleCredentials,
		in,
		false))
}
