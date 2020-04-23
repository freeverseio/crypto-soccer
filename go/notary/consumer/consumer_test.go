package consumer_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/consumer"
	"gotest.tools/assert"

	marketpay "github.com/freeverseio/crypto-soccer/go/marketpay/v1"
)

func TestConsumerNew(t *testing.T) {
	ch := make(chan interface{}, 10)
	_, err := consumer.New(
		ch,
		marketpay.NewMockMarketPay(),
		db,
		*bc.Contracts,
		bc.Owner,
	)
	assert.NilError(t, err)
}
