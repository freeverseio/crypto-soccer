package storage

import log "github.com/sirupsen/logrus"

type SellOffer struct {
	PlayerId uint64
	Price    uint64
}

func (b *Storage) CreateSellOffer(offer SellOffer) error {
	log.Infof("[DBMS] + create sell offer %v", offer)
	_, err := b.db.Exec("INSERT INTO player_sell_orders (playerId, price) VALUES ($1, $2);",
		offer.PlayerId,
		offer.Price,
	)
	return err
}

func (b *Storage) GetSellOfferts() ([]SellOffer, error) {
	var offers []SellOffer
	rows, err := b.db.Query("SELECT playerId, price FROM player_sell_orders;")
	if err != nil {
		return offers, err
	}
	defer rows.Close()
	for rows.Next() {
		var offer SellOffer
		err = rows.Scan(
			&offer.PlayerId,
			&offer.Price,
		)
		if err != nil {
			return offers, err
		}
		offers = append(offers, offer)
	}
	return offers, nil
}

func (b *Storage) DeleteSellOffer(playerId uint64) error {
	log.Infof("[DBMS] Delete sell offer %v", playerId)
	_, err := b.db.Exec("DELETE FROM player_sell_orders WHERE (playerId == '$1');",
		playerId,
	)
	return err
}
