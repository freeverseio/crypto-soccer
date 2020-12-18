package postgres

import (
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	log "github.com/sirupsen/logrus"
)

func (b *Tx) Unpayments(owner string) ([]*storage.Unpayment, error) {
	rows, err := b.tx.Query("SELECT id, time_of_unpayment, auction_id, notified FROM unpayment WHERE owner = $1;", owner)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var unpayments []*storage.Unpayment
	for rows.Next() {
		var unpayment storage.Unpayment
		unpayment.Owner = owner
		err = rows.Scan(
			&unpayment.Id,
			&unpayment.TimeOfUnpayment,
			&unpayment.AuctionId,
			&unpayment.Notified,
		)
		unpayments = append(unpayments, &unpayment)
	}
	return unpayments, err
}

func (b *Tx) UnpaymentInsert(unpayment storage.Unpayment) error {
	log.Debugf("[DBMS] + create unpayment %v", unpayment)
	_, err := b.tx.Exec(`INSERT INTO unpayment (owner, auction_id, time_of_unpayment) VALUES ($1, $2, NOW()) ON CONFLICT (auction_id, owner) DO NOTHING;`,
		unpayment.Owner,
		unpayment.AuctionId,
	)
	return err
}

func (b *Tx) UnpaymentUpdateNotified(unpayment storage.Unpayment) error {
	log.Debugf("[DBMS] + update unpayment %v", unpayment)
	_, err := b.tx.Exec(`UPDATE unpayment SET notified=$1 WHERE id=$2;`,
		unpayment.Notified,
		unpayment.Id,
	)
	return err
}
