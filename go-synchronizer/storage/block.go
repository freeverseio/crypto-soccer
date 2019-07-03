package storage

func (b *Storage) GetBlockNumber() (uint64, error) {
	rows, err := b.db.Query("SELECT value FROM params WHERE (name == 'block_number');")
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	if !rows.Next() {
		return 0, nil
	}
	var number uint64
	rows.Scan(&number)
	return number, nil
}

func (b *Storage) SetBlockNumber(value uint64) error {
	_, err := b.db.Exec("UPDATE params SET value = $1 WHERE (name == 'block_number');", value)
	if err != nil {
		return err
	}
	return nil
}
