package consumer_test

import (
	"testing"

	v1 "github.com/freeverseio/crypto-soccer/go/marketpay/v1"
	"github.com/freeverseio/crypto-soccer/go/notary/storage/postgres"

	"github.com/freeverseio/crypto-soccer/go/notary/consumer"
	"gotest.tools/assert"
)

func TestConsumerNew(t *testing.T) {
	ch := make(chan interface{}, 10)
	_, err := consumer.New(
		ch,
		v1.NewMockMarketPay(),
		*bc.Contracts,
		bc.Owner,
		googleCredentials,
		namesdb,
		false,
		postgres.NewStorageService(db),
	)
	assert.NilError(t, err)
	in := struct{}{}
	assert.Error(t, c.Consume(in), "unknown event: {}")
}

func TestConsumerConsumeCreateOffer(t *testing.T) {
	ch := make(chan interface{}, 10)
	c, err := consumer.New(
		ch,
		v1.NewMockMarketPay(),
		*bc.Contracts,
		bc.Owner,
		googleCredentials,
		namesdb,
		false,
		postgres.NewStorageService(db),
	)
	assert.NilError(t, err)
	in := input.CreateOfferInput{}
	assert.Error(t, c.Consume(in), "invalid teamId")
}
