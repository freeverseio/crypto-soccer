package playstore

import (
	"context"
	"crypto/ecdsa"
	"fmt"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/names"
	"github.com/freeverseio/crypto-soccer/go/notary/googleplaystoreutils"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	log "github.com/sirupsen/logrus"
)

type Machine struct {
	client    googleplaystoreutils.ClientService
	order     storage.PlaystoreOrder
	contracts contracts.Contracts
	pvc       *ecdsa.PrivateKey
	namesdb   *names.Generator
	iapTestOn bool
}

func New(
	client googleplaystoreutils.ClientService,
	order storage.PlaystoreOrder,
	contracts contracts.Contracts,
	pvc *ecdsa.PrivateKey,
	namesdb *names.Generator,
	iapTestOn bool,
) (*Machine, error) {
	return &Machine{
		client:    client,
		order:     order,
		contracts: contracts,
		pvc:       pvc,
		iapTestOn: iapTestOn,
		namesdb:   namesdb,
	}, nil
}

func (b Machine) Order() storage.PlaystoreOrder {
	return b.order
}

func (b *Machine) Process() error {
	ctx := context.Background()

	switch b.order.State {
	case storage.PlaystoreOrderOpen:
		log.Infof("[playstore|process] orderId %v in Open state", b.order.OrderId)
		return b.processOpenState(ctx)
	case storage.PlaystoreOrderAcknowledged:
		log.Infof("[playstore|process] orderId %v in Acknowledged state", b.order.OrderId)
		return b.processAcknowledged(ctx)
	case storage.PlaystoreOrderRefunding:
		log.Infof("[playstore|process] orderId %v in Refunding state", b.order.OrderId)
		return b.processRefundingState(ctx)
	case storage.PlaystoreOrderFailed:
		log.Warning("failed order ... skip")
		return nil
	case storage.PlaystoreOrderRefunded:
		log.Warning("refunded order ... skip")
		return nil
	case storage.PlaystoreOrderComplete:
		log.Warning("complete order ... skip")
		return nil
	default:
		return fmt.Errorf("unknown state %v", b.order.State)
	}
}

func (b *Machine) setState(state storage.PlaystoreOrderState, extra string) {
	if state == storage.PlaystoreOrderFailed {
		log.Warnf("order %v in state %v with %v", b.order.OrderId, state, extra)
	}
	b.order.State = state
	b.order.StateExtra = extra
}
