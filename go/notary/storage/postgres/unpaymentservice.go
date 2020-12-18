package postgres

import (
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	log "github.com/sirupsen/logrus"
)

func (b *Tx) Unpayments(owner string) ([]storage.Unpayment, error) {
	rows, err := b.tx.Query("SELECT time_of_unpayment, notified FROM unpayment WHERE owner = $1;", owner)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var unpayments []storage.Unpayment
	for rows.Next() {
		var unpayment storage.Unpayment
		unpayment.Owner = owner
		err = rows.Scan(
			&unpayment.TimeOfUnpayment,
			&unpayment.Notified,
		)
		unpayments = append(unpayments, unpayment)
	}
	return unpayments, err
}

func (b *Tx) UnpaymentInsert(unpayment storage.Unpayment) error {
	log.Debugf("[DBMS] + create unpayment %v", b)
	_, err := b.tx.Exec(`INSERT INTO unpayment (owner, time_of_unpayment) VALUES ($1, NOW());`,
		unpayment.Owner,
	)
	return err
}

func (b *Tx) UnpaymentUpdateNotified(unpayment storage.Unpayment) error {
	log.Debugf("[DBMS] + update unpayment %v", b)
	_, err := b.tx.Exec(`UPDATE unpayment SET notified=$1 WHERE id=$2;`,
		unpayment.Notified,
		unpayment.Id,
	)
	return err
}
