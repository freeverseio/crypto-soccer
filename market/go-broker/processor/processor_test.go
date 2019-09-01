package processor_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/market/go-broker/processor"
)

func TestProcess(t *testing.T) {
	processor := processor.NewProcessor()
	processor.Process()
}
