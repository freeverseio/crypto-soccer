package purchasevoider_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/purchasevoider"
	"github.com/freeverseio/crypto-soccer/go/purchasevoider/mockup"
	"gotest.tools/assert"
)

func TestNew(t *testing.T) {
	t.Run("nil param", func(t *testing.T) {
		_, err := purchasevoider.New(nil, nil, nil)
		assert.Error(t, err, "invalid params")
	})
	t.Run("new", func(t *testing.T) {
		_, err := purchasevoider.New(
			&mockup.VoidPurchaseService{},
			&mockup.UniverseService{},
			&mockup.MarketService{},
		)
		assert.NilError(t, err)
	})
}
