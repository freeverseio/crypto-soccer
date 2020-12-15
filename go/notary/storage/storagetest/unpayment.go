package storagetest

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"gotest.tools/assert"
)

func testUnpaymentServiceInterface(t *testing.T, service storage.StorageService) {
	t.Run("TestUnpaymentByOwnerUnexistent", func(t *testing.T) {
		tx, err := service.Begin()
		assert.NilError(t, err)
		defer tx.Rollback()

		unpayment, err := tx.Unpayment("4343")
		assert.NilError(t, err)
		assert.Assert(t, unpayment == nil)
	})

	t.Run("TestUnpaymentUpsert", func(t *testing.T) {
		tx, err := service.Begin()
		assert.NilError(t, err)
		defer tx.Rollback()

		unpayment := storage.NewUnpayment()
		unpayment.Owner = "ciao"
		unpayment.NumOfUnpayments = 0
		unpayment.LastTimeOfUnpayment = "3"
		assert.NilError(t, tx.UnpaymentUpsert(*unpayment))

		result, err := tx.Unpayment(unpayment.Owner)
		assert.NilError(t, err)
		assert.Equal(t, *result, *unpayment)

		assert.NilError(t, tx.UnpaymentUpsert(*unpayment))

		unpayment.NumOfUnpayments = 1
		result, err = tx.Unpayment(unpayment.Owner)
		assert.NilError(t, err)
		assert.Equal(t, *result, *unpayment)

	})

}
