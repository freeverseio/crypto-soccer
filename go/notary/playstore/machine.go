package playstore

import (
	"fmt"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	log "github.com/sirupsen/logrus"
)

type Machine struct {
	googleCredentials []byte
	order             storage.PlaystoreOrder
}

func New(
	googleCredentials []byte,
	order storage.PlaystoreOrder,
) *Machine {
	return &Machine{
		googleCredentials: googleCredentials,
		order:             order,
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
