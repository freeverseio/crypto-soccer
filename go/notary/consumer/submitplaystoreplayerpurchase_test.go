package consumer_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/consumer"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/storage/postgres"
	"gotest.tools/assert"
)

func TestSubmitPlayStorePlayerPurchaseEmptyEvent(t *testing.T) {
	service := postgres.NewStorageService(db)
	tx, err := service.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	in := input.SubmitPlayStorePlayerPurchaseInput{}
	assert.Error(t, consumer.SubmitPlayStorePlayerPurchase(tx, in), "unexpected end of JSON input")
}

func TestSubmitPlayStorePlayerPurchase(t *testing.T) {
	service := postgres.NewStorageService(db)
	tx, err := service.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	in := input.SubmitPlayStorePlayerPurchaseInput{}
	in.PlayerId = "3"
	in.TeamId = "4"
	in.Receipt = "{\"Store\":\"GooglePlay\",\"TransactionID\":\"GPA.3306-7253-6203-28834\",\"Payload\":\"{\\\"json\\\":\\\"{\\\\\\\"orderId\\\\\\\":\\\\\\\"GPA.3306-7253-6203-28834\\\\\\\",\\\\\\\"packageName\\\\\\\":\\\\\\\"com.freeverse.phoenix\\\\\\\",\\\\\\\"productId\\\\\\\":\\\\\\\"player_tier_0\\\\\\\",\\\\\\\"purchaseTime\\\\\\\":1589875478821,\\\\\\\"purchaseState\\\\\\\":0,\\\\\\\"developerPayload\\\\\\\":\\\\\\\"{\\\\\\\\\\\\\\\"developerPayload\\\\\\\\\\\\\\\":\\\\\\\\\\\\\\\"\\\\\\\\\\\\\\\",\\\\\\\\\\\\\\\"is_free_trial\\\\\\\\\\\\\\\":false,\\\\\\\\\\\\\\\"has_introductory_price_trial\\\\\\\\\\\\\\\":false,\\\\\\\\\\\\\\\"is_updated\\\\\\\\\\\\\\\":false,\\\\\\\\\\\\\\\"accountId\\\\\\\\\\\\\\\":\\\\\\\\\\\\\\\"\\\\\\\\\\\\\\\"}\\\\\\\",\\\\\\\"purchaseToken\\\\\\\":\\\\\\\"jgflonflnpkpjnnlcjgbbenh.AO-J1OyzOeK1EnYJr5vVpD6fEc6T0IQZKTFAeE80Z1Q7XXMdwBxYvu3cP0HlYWv4lQ7lH5gllAz-YnWMjouOoc011JE_rPtPQlzLy5sA4sv-Lo8apTomkV20POaAFilMnt2GSZOnFeeh\\\\\\\"}\\\",\\\"signature\\\":\\\"EUbadmcvm0IqVMsedO79UvQWFfzwBr2OimQfKI5Md76UQiDxPr1sSSGCYz4ln807ryQLBbG\\\\/PjBIAjBuFRyArS5DgZN3ngtetGebzuJ9plYfAj+NXbAStErchp95rpmpW+Z5DCe0O7DqFYstsWGhLmswA5uujUEDzNvy0WZRddy1LudyioN1wi3VlOuGSDdOWq2pV91Gx4xGlSANhSlafrvIuBFQVNBDQFSc9YhTcrB5iEqcURh8V6kKizuqxMeSWKsLqrU8AxMQfclKc1w+EWOfHerSAA1hWroxl4wT165UhXwNtoK8XLR7J+Ymc0nvFUJRY0fYYEx0dc2WjPPfdg==\\\",\\\"skuDetails\\\":\\\"{\\\\\\\"skuDetailsToken\\\\\\\":\\\\\\\"AEuhp4LPTDFb7nqP3iZ-YPCreWxMx1tBtAHRTWI5pucDd9R8TCw8cf4NlsYOEH0Ix5wr\\\\\\\",\\\\\\\"productId\\\\\\\":\\\\\\\"player_tier_0\\\\\\\",\\\\\\\"type\\\\\\\":\\\\\\\"inapp\\\\\\\",\\\\\\\"price\\\\\\\":\\\\\\\"0,59\u00a0€\\\\\\\",\\\\\\\"price_amount_micros\\\\\\\":590000,\\\\\\\"price_currency_code\\\\\\\":\\\\\\\"EUR\\\\\\\",\\\\\\\"title\\\\\\\":\\\\\\\"Player Tier 0 (goalRevolution)\\\\\\\",\\\\\\\"description\\\\\\\":\\\\\\\"This player isn't good, but he's cheap.\\\\\\\"}\\\",\\\"isPurchaseHistorySupported\\\":true}\"}"
	assert.NilError(t, consumer.SubmitPlayStorePlayerPurchase(tx, in))
}
