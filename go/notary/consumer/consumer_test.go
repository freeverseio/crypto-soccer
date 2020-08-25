package consumer_test

import (
	"testing"

	v1 "github.com/freeverseio/crypto-soccer/go/marketpay/v1"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"

	"github.com/freeverseio/crypto-soccer/go/notary/consumer"
	"gotest.tools/assert"
)

func TestConsumerNew(t *testing.T) {
	ch := make(chan interface{}, 10)
	_, err := consumer.New(
		ch,
		v1.NewMockMarketPay(),
		db,
		*bc.Contracts,
		bc.Owner,
		googleCredentials,
		namesdb,
		false,
		service,
	)
	assert.NilError(t, err)
}

func TestConsumerConsumeSubmitPlayStorePlayerPurchaseInput(t *testing.T) {
	ch := make(chan interface{}, 10)
	c, err := consumer.New(
		ch,
		v1.NewMockMarketPay(),
		db,
		*bc.Contracts,
		bc.Owner,
		googleCredentials,
		namesdb,
		false,
		service,
	)
	assert.NilError(t, err)
	in := input.SubmitPlayStorePlayerPurchaseInput{}
	assert.Error(t, c.Consume(in), "unexpected end of JSON input")
}

func TestConsumerConsumeCreateAuction(t *testing.T) {
	ch := make(chan interface{}, 10)
	c, err := consumer.New(
		ch,
		v1.NewMockMarketPay(),
		db,
		*bc.Contracts,
		bc.Owner,
		googleCredentials,
		namesdb,
		false,
		service,
	)
	assert.NilError(t, err)
	in := input.CreateAuctionInput{}
	assert.Error(t, c.Consume(in), "invalid playerId")
}

func TestConsumerConsumeCancelAuction(t *testing.T) {
	ch := make(chan interface{}, 10)
	c, err := consumer.New(
		ch,
		v1.NewMockMarketPay(),
		db,
		*bc.Contracts,
		bc.Owner,
		googleCredentials,
		namesdb,
		false,
		service,
	)
	assert.NilError(t, err)
	in := input.CancelAuctionInput{}
	assert.Error(t, c.Consume(in), "trying to cancel a nil auction")
}

func TestConsumerConsumeCreateBid(t *testing.T) {
	ch := make(chan interface{}, 10)
	c, err := consumer.New(
		ch,
		v1.NewMockMarketPay(),
		db,
		*bc.Contracts,
		bc.Owner,
		googleCredentials,
		namesdb,
		false,
		service,
	)
	assert.NilError(t, err)
	in := input.CreateBidInput{}
	assert.Error(t, c.Consume(in), "No auction for bid {  0 0 }")
}

func TestConsumerConsumeUnknownEvent(t *testing.T) {
	ch := make(chan interface{}, 10)
	c, err := consumer.New(
		ch,
		v1.NewMockMarketPay(),
		db,
		*bc.Contracts,
		bc.Owner,
		googleCredentials,
		namesdb,
		false,
		service,
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
		db,
		*bc.Contracts,
		bc.Owner,
		googleCredentials,
		namesdb,
		false,
		service,
	)
	assert.NilError(t, err)
	in := input.CreateOfferInput{}
	assert.Error(t, c.Consume(in), "invalid teamId")
}
