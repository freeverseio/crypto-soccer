package testutils_test

import (
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
	bc.DeployContracts(bc.Owner)
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
}
