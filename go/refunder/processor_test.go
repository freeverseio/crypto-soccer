package refunder_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/refunder"
	"github.com/freeverseio/crypto-soccer/go/refunder/mockup"
	"gotest.tools/assert"
)

func TestNew(t *testing.T) {
	t.Run("nil param", func(t *testing.T) {
		_, err := refunder.New(nil, nil, nil)
		assert.Error(t, err, "invalid params")
	})
	t.Run("new", func(t *testing.T) {
		_, err := refunder.New(
			&mockup.PaymentService{},
			&mockup.UniverseService{},
			&mockup.MarketService{},
		)
		assert.NilError(t, err)
	})
}
