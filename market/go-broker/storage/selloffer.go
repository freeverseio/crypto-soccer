package storage

type SellOffer struct {
	PlayerId uint64
	Price    uint64
}

func (b *Storage) TeamCount() (uint64, error) {
	rows, err := b.db.Query("SELECT COUNT(*) FROM teams;")
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	rows.Next()
	var count uint64
	err = rows.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
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
