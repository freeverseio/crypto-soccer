package postgres

import (
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

func (b *Tx) AuctionPass(owner string) (*storage.AuctionPass, error) {
	rows, err := b.tx.Query(`SELECT 
	owner,
	purchased_for_team_id,
	product_id
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
		&order.PurchasedForTeamId,
		&order.ProductId,
	)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (b *Tx) AuctionPassInsert(order storage.AuctionPass) error {
	_, err := b.tx.Exec(`INSERT INTO auction_pass (
		owner,
		purchased_for_team_id,
		product_id
		) VALUES ($1);`,
		order.Owner,
		order.PurchasedForTeamId,
		order.ProductId,
	)
	return err
}

func (b *Tx) AuctionPassAcknowledge(ap storage.AuctionPass) error {
	_, err := b.tx.Exec(`UPDATE auction_pass SET 
		ack=$1, 
		WHERE owner=$2;`,
		ap.Ack,
		ap.Owner,
	)
	return err
}
