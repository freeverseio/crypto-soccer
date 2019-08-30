package storage

type SellOffer struct {
	PlayerId uint64
	Price    uint64
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
