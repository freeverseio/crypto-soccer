package consumer_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/consumer"
	"gotest.tools/assert"
)

func TestConsumerNew(t *testing.T) {
	ch := make(chan interface{}, 10)
	_, err := consumer.New(
		ch,
		db,
		*bc.Contracts,
		bc.Owner,
	)
	assert.NilError(t, err)
}
