package postgres

import (
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	log "github.com/sirupsen/logrus"
)

func (b *Tx) Unpayment(owner string) (*storage.Unpayment, error) {
	rows, err := b.tx.Query("SELECT num_of_unpayments, last_time_of_unpayment FROM unpayment WHERE owner = $1;", owner)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, nil
	}
	var unpayment storage.Unpayment
	unpayment.Owner = owner
	err = rows.Scan(
		&unpayment.NumOfUnpayments,
		&unpayment.LastTimeOfUnpayment,
	)
	return &unpayment, err
}

func (b *Tx) UnpaymentUpsert(unpayment storage.Unpayment) error {
	log.Debugf("[DBMS] + create or update unpayment %v", b)
	_, err := b.tx.Exec(`INSERT INTO unpayment (owner, num_of_unpayments, last_time_of_unpayment) VALUES ($1, $2, NOW()) 
	ON CONFLICT(owner) DO UPDATE SET num_of_unpayments=num_of_unpayments + 1, last_time_of_unpayment=NOW();`,
		unpayment.Owner,
		unpayment.NumOfUnpayments,
		unpayment.LastTimeOfUnpayment,
	)
	return err
}
