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

		unpayment, err := tx.Unpayments("4343")
		assert.NilError(t, err)
		assert.Assert(t, unpayment == nil)
	})

	t.Run("TestUnpaymentInsert", func(t *testing.T) {
		tx, err := service.Begin()
		assert.NilError(t, err)
		defer tx.Rollback()

		unpayment := storage.NewUnpayment()
		unpayment.Owner = "ciao"
		unpayment.LastTimeOfUnpayment = "3"
		assert.NilError(t, tx.UnpaymentInsert(*unpayment))

		result, err := tx.Unpayments(unpayment.Owner)
		assert.NilError(t, err)
		assert.Equal(t, len(result), 1)
		assert.Equal(t, result[0], unpayment)

		assert.NilError(t, tx.UnpaymentInsert(*unpayment))

		result, err = tx.Unpayments(unpayment.Owner)
		assert.NilError(t, err)
		assert.Equal(t, len(result), 2)

	})

}
