package postgres

import (
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

func (b *Tx) AuctionPass(owner string) (*storage.AuctionPass, error) {
	rows, err := b.tx.Query(`SELECT 
	owner
	FROM auction_pass WHERE owner=$1;`, owner)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	order := storage.AuctionPass{}

	err = rows.Scan(
		&order.Owner,
	)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (b *Tx) AuctionPassInsert(order storage.AuctionPass) error {
	_, err := b.tx.Exec(`INSERT INTO auction_pass (
		owner
		) VALUES ($1);`,
		order.Owner,
	)
	return err
}
