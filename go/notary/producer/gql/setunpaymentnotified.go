package gql

import (
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	log "github.com/sirupsen/logrus"
)

func (b *Resolver) SetUnpaymentNotified(args struct {
	Input input.SetUnpaymentNotifiedInput
}) (*bool, error) {
	log.Debugf("SetUnpaymentNotified %v", args)

	tx, err := b.service.Begin()
	owner := string(args.Input.Owner)
	unpayment := storage.NewUnpayment()
	unpayment.Owner = owner
	unpayment.Notified = true
	err = tx.UnpaymentUpdateNotified(*unpayment)
	ok := true
	if err != nil {
		ok = false
		tx.Rollback()
		return &ok, err
	}

	return &ok, tx.Commit()
}
