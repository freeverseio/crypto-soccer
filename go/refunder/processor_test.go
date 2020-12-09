package refunder_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/refunder"
	mmockup "github.com/freeverseio/crypto-soccer/go/refunder/storage/market/mockup"
	umockup "github.com/freeverseio/crypto-soccer/go/refunder/storage/universe/mockup"
	"gotest.tools/assert"
)

func TestNew(t *testing.T) {
	t.Run("nil param", func(t *testing.T) {
		_, err := refunder.New(nil, nil)
		assert.Error(t, err, "invalid params")
	})
	t.Run("new", func(t *testing.T) {
		_, err := refunder.New(&umockup.Service{}, &mmockup.Service{})
		assert.NilError(t, err)
	})
}
