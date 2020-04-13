package consumer_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/consumer"
	"gotest.tools/assert"
)

func TestConsumerNew(t *testing.T) {
	ch := make(chan interface{}, 100000)
	_, err := consumer.New(
		ch,
		db,
	)
	assert.NilError(t, err)
}
