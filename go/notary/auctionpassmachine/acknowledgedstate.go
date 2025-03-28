package auctionpassmachine

import (
	"context"
	"errors"
	"math/big"

	"github.com/freeverseio/crypto-soccer/go/notary/googleplaystoreutils"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"

	log "github.com/sirupsen/logrus"
)

func (b *AuctionPassMachine) processAuctionPassAcknowledged(ctx context.Context, service storage.Tx) error {
	purchase, err := b.client.GetPurchase(
		ctx,
		b.order.PackageName,
		b.order.ProductId,
		b.order.PurchaseToken,
	)
	if err != nil {
		b.setState(storage.AuctionPassPlaystoreOrderAcknowledged, err.Error())
		return nil
	}

	validator := googleplaystoreutils.NewPurchaseValidator(*purchase)
	if validator.IsTest() && !b.iapTestOn {
		log.Warningf("[consumer|iap] received test orderId %v ...", purchase.OrderId)
		b.setState(storage.AuctionPassPlaystoreOrderComplete, "test order")
	} else {
		if err := b.acknowledgeAuctionPass(service); err != nil {
			b.setState(storage.AuctionPassPlaystoreOrderRefunding, err.Error())
			return nil
		}
		log.Infof("[consumer|iap] auction pass orderId %v owner %v assigned to teamId %v", purchase.OrderId, b.order.Owner, b.order.TeamId)
		b.setState(storage.AuctionPassPlaystoreOrderComplete, "")
	}

	return nil
}

func (b AuctionPassMachine) acknowledgeAuctionPass(service storage.Tx) error {

	teamId, _ := new(big.Int).SetString(b.order.TeamId, 10)
	if teamId == nil {
		return errors.New("invalid team")
	}

	auctionPass := storage.AuctionPass{
		Owner: string(b.order.Owner),
	}
	err := service.AuctionPassInsert(auctionPass)
	if err != nil {
		return err
	}

	return nil
}
