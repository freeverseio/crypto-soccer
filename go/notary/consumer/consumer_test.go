package consumer_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/consumer"
)

func TestConsumerNew(t *testing.T) {
	ch := make(chan interface{}, 100000)
	c := consumer.New(
		ch,
	)
}
