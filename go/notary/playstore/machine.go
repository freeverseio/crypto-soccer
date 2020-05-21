package playstore

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	log "github.com/sirupsen/logrus"
)

type Machine struct {
	googleCredentials []byte
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
) *Machine {

	return &Machine{
		googleCredentials: googleCredentials,
		order:             order,
		contracts:         contracts,
		pvc:               pvc,
		iapTestOn:         iapTestOn,
	}
}

func (b Machine) Order() storage.PlaystoreOrder {
	return b.order
}

func (b *Machine) Process() error {
	switch b.order.State {
	case storage.PlaystoreOrderPending:
		return b.processPendingState()
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
