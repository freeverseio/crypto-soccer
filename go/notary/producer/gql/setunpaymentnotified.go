package gql

import (
	"errors"
	"strconv"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	log "github.com/sirupsen/logrus"
)

func (b *Resolver) SetUnpaymentNotified(args struct {
	Input input.SetUnpaymentNotifiedInput
}) (*bool, error) {
	log.Debugf("SetUnpaymentNotified %v", args)

	ok := true
	tx, err := b.service.Begin()
	id, err := strconv.ParseInt(string(args.Input.Id), 10, 64)
	if err != nil {
		ok = false
		return &ok, errors.New("Invalid id")
	}
	unpayment := storage.NewUnpayment()
	unpayment.Id = id
	unpayment.Notified = true
	err = tx.UnpaymentUpdateNotified(*unpayment)
	if err != nil {
		ok = false
		tx.Rollback()
		return &ok, err
	}

	return &ok, tx.Commit()
}
