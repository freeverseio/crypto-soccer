package playstore

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/helper"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

func (b *Machine) processOpenState(ctx context.Context) error {
	purchase, err := b.client.VerifyProduct(
		ctx,
		b.order.PackageName,
		b.order.ProductId,
		b.order.PurchaseToken,
	)
	if err != nil {
		b.setState(storage.PlaystoreOrderFailed, err.Error())
		return nil
	}

	validator := NewPurchaseValidator(*purchase)
	if validator.IsCanceled() {
		b.setState(storage.PlaystoreOrderComplete, "cancelled")
		return nil
	}
	if validator.IsPending() {
		b.setState(storage.PlaystoreOrderOpen, "pending")
		return nil
	}
	if !validator.IsPurchased() {
		b.setState(storage.PlaystoreOrderFailed, "invalid puchase state")
		return nil
	}
	if validator.IsAcknowledged() {
		b.setState(storage.PlaystoreOrderFailed, "already acknowledged")
		return nil
	}

	payload := fmt.Sprintf("playerId: %v", b.order.PlayerId)
	if err := b.client.AcknowledgeProduct(
		ctx,
		b.order.PackageName,
		b.order.ProductId,
		b.order.PurchaseToken,
		payload,
	); err != nil {
		b.setState(storage.PlaystoreOrderFailed, err.Error())
		return err
	}

	b.setState(storage.PlaystoreOrderAcknowledged, "")
	return nil
}

func (b Machine) assignAsset() error {
	playerId, _ := new(big.Int).SetString(b.order.PlayerId, 10)
	if playerId == nil {
		return errors.New("invalid player")
	}
	teamId, _ := new(big.Int).SetString(b.order.TeamId, 10)
	if teamId == nil {
		return errors.New("invalid team")
	}

	auth := bind.NewKeyedTransactor(b.pvc)
	auth.GasPrice = big.NewInt(1000000000) // in xdai is fixe to 1 GWei
	tx, err := b.contracts.Market.TransferBuyNowPlayer(
		auth,
		playerId,
		teamId,
	)
	if err != nil {
		return err
	}
	if _, err = helper.WaitReceipt(b.contracts.Client, tx, 60); err != nil {
		return err
	}
	return nil
}
