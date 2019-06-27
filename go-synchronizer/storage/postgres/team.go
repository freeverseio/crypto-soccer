package postgres

import "github.com/freeverseio/crypto-soccer/go-synchronizer/storage"

func (b *PostgresStorage) TeamAdd(ID uint64, name string) error {
	//  TODO: check for db is initialized
	_, err := b.db.Query("INSERT INTO teams (id, name) VALUES ($1, $2);", ID, name)
	if err != nil {
		return err
	}

	return nil
}

func (b *PostgresStorage) TeamCount() (uint64, error) {
	var count uint64
	row := b.db.QueryRow("SELECT COUNT(*) FROM teams;")
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (b *PostgresStorage) GetTeam(id uint64) (storage.Team, error) {
	return storage.Team{}, nil
}
