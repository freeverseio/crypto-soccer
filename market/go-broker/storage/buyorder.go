package storage

type BuyOffer struct {
	PlayerId uint64
	Price    uint64
}

func (b *Storage) GetBuyOfferts() ([]BuyOffer, error) {
	var offers []BuyOffer
	rows, err := b.db.Query("SELECT playerId, price FROM player_buy_orders;")
	if err != nil {
		return offers, err
	}
	defer rows.Close()
	for rows.Next() {
		var offer BuyOffer
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
