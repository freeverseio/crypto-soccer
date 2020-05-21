package playstore

import (
	"context"
	"crypto/ecdsa"
	"fmt"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	log "github.com/sirupsen/logrus"
)

type Machine struct {
	client    ClientService
	order     storage.PlaystoreOrder
	contracts contracts.Contracts
	pvc       *ecdsa.PrivateKey
	iapTestOn bool
}

func New(
	client ClientService,
	order storage.PlaystoreOrder,
	contracts contracts.Contracts,
	pvc *ecdsa.PrivateKey,
	iapTestOn bool,
) (*Machine, error) {
	return &Machine{
		client:    client,
		order:     order,
		contracts: contracts,
		pvc:       pvc,
		iapTestOn: iapTestOn,
	}, nil
}

func (b Machine) Order() storage.PlaystoreOrder {
	return b.order
}

func (b *Machine) Process() error {
	ctx := context.Background()

	switch b.order.State {
	case storage.PlaystoreOrderOpen:
		return b.processOpenState(ctx)
	case storage.PlaystoreOrderAcknowledged:
		return b.processAcknowledged(ctx)
	case storage.PlaystoreOrderRefunding:
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
