package playstore

import (
	"context"
	"crypto/ecdsa"
	"fmt"

	"github.com/awa/go-iap/playstore"
	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	log "github.com/sirupsen/logrus"
)

type Machine struct {
	googleCredentials []byte
	client            *playstore.Client
	order             storage.PlaystoreOrder
	contracts         contracts.Contracts
	pvc               *ecdsa.PrivateKey
	iapTestOn         bool
}

func New(
	googleCredentials []byte,
	order storage.PlaystoreOrder,
	contracts contracts.Contracts,
	pvc *ecdsa.PrivateKey,
	iapTestOn bool,
) (*Machine, error) {
	client, err := playstore.New(googleCredentials)
	if err != nil {
		return nil, err
	}

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
		return b.processPendingState(ctx)
	case storage.PlaystoreOrderAssetAssigned:
		return b.processAssetAssigned(ctx)
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
