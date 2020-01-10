package storage

import "database/sql"

func GetBlockNumber(tx *sql.Tx) (uint64, error) {
	rows, err := tx.Query("SELECT value FROM params WHERE (name = 'block_number');")
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

func SetBlockNumber(tx *sql.Tx, value uint64) error {
	_, err := tx.Exec("UPDATE params SET value = $1 WHERE (name = 'block_number');", value)
	if err != nil {
		return err
	}
	return nil
}
