package purchasevoider_test

import (
	"context"
	"errors"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/purchasevoider"
	"github.com/freeverseio/crypto-soccer/go/purchasevoider/mockup"
	"google.golang.org/api/androidpublisher/v3"
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

func TestGetVoidedTokens(t *testing.T) {
	t.Run("service return nil tockens", func(t *testing.T) {
		pv, err := purchasevoider.New(
			&mockup.VoidPurchaseService{
				VoidedPurchasesFn: func(context.Context) ([]*androidpublisher.VoidedPurchase, error) {
					return nil, nil
				},
			},
			&mockup.UniverseService{},
			&mockup.MarketService{},
		)
		assert.NilError(t, err)
		tokens, err := pv.GetVoidedTokens()
		assert.NilError(t, err)
		assert.Equal(t, len(tokens), 0)
	})
	t.Run("service return error", func(t *testing.T) {
		pv, err := purchasevoider.New(
			&mockup.VoidPurchaseService{
				VoidedPurchasesFn: func(context.Context) ([]*androidpublisher.VoidedPurchase, error) {
					return nil, errors.New("error")
				},
			},
			&mockup.UniverseService{},
			&mockup.MarketService{},
		)
		assert.NilError(t, err)
		_, err = pv.GetVoidedTokens()
		assert.Error(t, err, "error")
	})
	t.Run("service return tokens", func(t *testing.T) {
		pv, err := purchasevoider.New(
			&mockup.VoidPurchaseService{
				VoidedPurchasesFn: func(context.Context) ([]*androidpublisher.VoidedPurchase, error) {
					return []*androidpublisher.VoidedPurchase{
						&androidpublisher.VoidedPurchase{
							PurchaseToken: "ciao",
						},
					}, nil
				},
			},
			&mockup.UniverseService{},
			&mockup.MarketService{},
		)
		assert.NilError(t, err)
		tokens, err := pv.GetVoidedTokens()
		assert.NilError(t, err)
		assert.Equal(t, len(tokens), 1)
		assert.Equal(t, tokens[0], "ciao")
	})
}

func TestGetPlayerIds(t *testing.T) {
	t.Run("0 tokens", func(t *testing.T) {
		pv, err := purchasevoider.New(
			&mockup.VoidPurchaseService{},
			&mockup.UniverseService{},
			&mockup.MarketService{},
		)
		assert.NilError(t, err)
		ids, err := pv.GetPlayerIds([]string{})
		assert.NilError(t, err)
		assert.Equal(t, len(ids), 0)
	})
	t.Run("no id from service", func(t *testing.T) {
		pv, err := purchasevoider.New(
			&mockup.VoidPurchaseService{},
			&mockup.UniverseService{},
			&mockup.MarketService{
				GetPlayerIdByPurchaseTokenFunc: func(token string) (string, error) {
					return "", nil
				},
			},
		)
		assert.NilError(t, err)
		ids, err := pv.GetPlayerIds([]string{"ciao"})
		assert.NilError(t, err)
		assert.Equal(t, len(ids), 0)
	})
	t.Run("error from service", func(t *testing.T) {
		pv, err := purchasevoider.New(
			&mockup.VoidPurchaseService{},
			&mockup.UniverseService{},
			&mockup.MarketService{
				GetPlayerIdByPurchaseTokenFunc: func(token string) (string, error) {
					return "", errors.New("error")
				},
			},
		)
		assert.NilError(t, err)
		_, err = pv.GetPlayerIds([]string{"ciao"})
		assert.Error(t, err, "error")
	})
	t.Run("service returns ids", func(t *testing.T) {
		pv, err := purchasevoider.New(
			&mockup.VoidPurchaseService{},
			&mockup.UniverseService{},
			&mockup.MarketService{
				GetPlayerIdByPurchaseTokenFunc: func(token string) (string, error) {
					return "id", nil
				},
			},
		)
		assert.NilError(t, err)
		ids, err := pv.GetPlayerIds([]string{"ciao"})
		assert.NilError(t, err)
		assert.Equal(t, len(ids), 1)
		assert.Equal(t, ids[0], "id")
	})
}
