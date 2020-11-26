package auctionpassmachine

import (
	"context"
	"crypto/ecdsa"
	"fmt"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	log "github.com/sirupsen/logrus"
)

type AuctionPassMachine struct {
	service   storage.Tx
	client    ClientService
	order     storage.AuctionPassPlaystoreOrder
	contracts contracts.Contracts
	pvc       *ecdsa.PrivateKey
	iapTestOn bool
}

func New(
	service storage.Tx,
	client ClientService,
	order storage.AuctionPassPlaystoreOrder,
	contracts contracts.Contracts,
	pvc *ecdsa.PrivateKey,
	iapTestOn bool,
) (*AuctionPassMachine, error) {
	return &AuctionPassMachine{
		service:   service,
		client:    client,
		order:     order,
		contracts: contracts,
		pvc:       pvc,
		iapTestOn: iapTestOn,
	}, nil
}

func (b AuctionPassMachine) Order() storage.AuctionPassPlaystoreOrder {
	return b.order
}

func (b *AuctionPassMachine) Process() error {
	ctx := context.Background()

	switch b.order.State {
	case storage.AuctionPassPlaystoreOrderOpen:
		log.Infof("[auctionpass|playstore|process] orderId %v in Open state", b.order.OrderId)
		return b.processAuctionPassOpenState(ctx, b.service)
	case storage.AuctionPassPlaystoreOrderAcknowledged:
		log.Infof("[auctionpass|playstore|process] orderId %v in Acknowledged state", b.order.OrderId)
		return b.processAuctionPassAcknowledged(ctx, b.service)
	case storage.AuctionPassPlaystoreOrderRefunding:
		log.Infof("[auctionpass|playstore|process] orderId %v in Refunding state", b.order.OrderId)
		return b.processAuctionPassRefundingState(ctx)
	case storage.AuctionPassPlaystoreOrderFailed:
		log.Warning("failed order ... skip")
		return nil
	case storage.AuctionPassPlaystoreOrderRefunded:
		log.Warning("refunded order ... skip")
		return nil
	case storage.AuctionPassPlaystoreOrderComplete:
		log.Warning("complete order ... skip")
		return nil
	default:
		return fmt.Errorf("unknown state %v", b.order.State)
	}
}

func (b *AuctionPassMachine) setState(state storage.AuctionPassPlaystoreOrderState, extra string) {
	if state == storage.AuctionPassPlaystoreOrderFailed {
		log.Warnf("order %v in state %v with %v", b.order.OrderId, state, extra)
	}
	b.order.State = state
	b.order.StateExtra = extra
}
