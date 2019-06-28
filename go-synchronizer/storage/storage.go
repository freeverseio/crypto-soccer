package storage

import (
	"database/sql"
	"errors"
	"io/ioutil"
	"log"
	"math/big"
	"os"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type Team struct {
	Id   uint64
	Name string
}

type Storage struct {
	db *sql.DB
}

func NewPostgres(url string) (*Storage, error) {
	var err error
	storage := &Storage{}
	storage.db, err = sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	return storage, nil
}

func NewSqlite3(schemaFile string) (*Storage, error) {
	var err error
	storage := Storage{}
	storage.db, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}
	if err := storage.db.Ping(); err != nil {
		log.Fatalf("could not ping DB... %v", err)
	}
	file, err := os.Open(schemaFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	script, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	_, err = storage.db.Exec(string(script))
	if err != nil {
		log.Fatal(err)
	}
	return &storage, nil
}

func (b *Storage) TeamAdd(ID uint64, name string) error {
	//  TODO: check for db is initialized
	_, err := b.db.Exec("INSERT INTO teams (id, name) VALUES ($1, $2);", ID, name)
	if err != nil {
		return err
	}

	return nil
}

func (b *Storage) TeamCount() (uint64, error) {
	rows, err := b.db.Query("SELECT COUNT(*) FROM teams;")
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	rows.Next()
	var count uint64
	rows.Scan(&count)
	return count, nil
}

func (b *Storage) GetTeam(id uint64) (Team, error) {
	team := Team{}
	rows, err := b.db.Query("SELECT id, name FROM teams WHERE (id == $1);", id)
	if err != nil {
		return team, err
	}
	defer rows.Close()
	if !rows.Next() {
		return team, errors.New("unexistent team")
	}
	rows.Scan(&team.Id, &team.Name)
	return team, nil
}

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
