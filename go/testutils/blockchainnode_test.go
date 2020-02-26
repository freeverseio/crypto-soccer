package testutils_test

import (
	"fmt"
	"log"
	"testing"

	"gotest.tools/assert"

	"github.com/freeverseio/crypto-soccer/go/testutils"
)

func TestConctactsDeploy(t *testing.T) {
	bc, err := testutils.NewBlockchainNode()
	if err != nil {
		log.Fatal(err)
	}
	assert.NilError(t, bc.DeployContracts(bc.Owner))
	t.Log(fmt.Sprintf("%+v", bc.Contracts))
	assert.Assert(t, bc.Contracts.Client != nil)
	assert.Assert(t, bc.Contracts.Leagues != nil)
	assert.Assert(t, bc.Contracts.Assets != nil)
	assert.Assert(t, bc.Contracts.Evolution != nil)
	assert.Assert(t, bc.Contracts.Engine != nil)
	assert.Assert(t, bc.Contracts.Engineprecomp != nil)
	assert.Assert(t, bc.Contracts.Updates != nil)
	assert.Assert(t, bc.Contracts.Market != nil)
	assert.Assert(t, bc.Contracts.Utils != nil)
	assert.Assert(t, bc.Contracts.PlayAndEvolve != nil)
	assert.Assert(t, bc.Contracts.Shop != nil)
	assert.Assert(t, bc.Contracts.TrainingPoints != nil)
	assert.Assert(t, bc.Contracts.ConstantsGetters != nil)
}
