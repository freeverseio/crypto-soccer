package truffle_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/contracts/truffle"
	"gotest.tools/assert"
)

func TestTruffleDeploy(t *testing.T) {
	proxyAddress, err := truffle.Deploy()
	assert.NilError(t, err)
	assert.Assert(t, proxyAddress != "")
}

func TestTruffleNew(t *testing.T) {
	contracts, err := truffle.New()
	assert.NilError(t, err)
	assert.Assert(t, contracts.ProxyAddress != "")
}
