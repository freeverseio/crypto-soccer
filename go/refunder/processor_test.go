package refunder_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/refunder"
	"github.com/freeverseio/crypto-soccer/go/refunder/storage/mockup"
	"gotest.tools/assert"
)

func TestNew(t *testing.T) {
	t.Run("nil param", func(t *testing.T) {
		_, err := refunder.New(nil, nil)
		assert.Error(t, err, "invalid params")
	})
	t.Run("new", func(t *testing.T) {
		_, err := refunder.New(&mockup.UniverseService{}, &mockup.MarketService{})
		assert.NilError(t, err)
	})
}
