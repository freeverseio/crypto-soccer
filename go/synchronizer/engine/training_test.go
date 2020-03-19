package engine_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/synchronizer/engine"
	"gotest.tools/golden"
)

func TestTrainingToString(t *testing.T) {
	training := engine.NewTraining()
	golden.Assert(t, training.ToString(), t.Name()+".golden")
}
