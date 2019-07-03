package storage

import "math/big"

func (b *Storage) GetBlockNumber() (*big.Int, error) {
	rows, err := b.db.Query("SELECT value FROM params WHERE (name == 'block_number');")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, nil
	}
	var number int64
	rows.Scan(&number)
	return big.NewInt(number), nil
}

func (b *Storage) SetBlockNumber(value *big.Int) error {
	_, err := b.db.Exec("UPDATE params SET value = $1 WHERE (name == 'block_number');", value.Uint64())
	if err != nil {
		return err
	}
	return nil
}
