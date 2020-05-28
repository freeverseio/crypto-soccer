package contracts_test

import (
	"testing"

	"gotest.tools/assert"
)

func TestTruffleNew(t *testing.T) {
	t.Parallel()
	assert.Assert(t, bc.ProxyAddress != "")
}
