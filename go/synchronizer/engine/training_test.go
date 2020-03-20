package engine_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/storage"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/engine"
	"gotest.tools/golden"
)

func TestTrainingToString(t *testing.T) {
	training := engine.NewTraining(*storage.NewTraining())
	golden.Assert(t, training.Marshal(), t.Name()+".golden")
}
