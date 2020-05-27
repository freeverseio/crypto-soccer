package postgres

import (
	"database/sql"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/storage"
)

func NewFromStorage(client *ethclient.Client, tx *sql.Tx) (*contracts.Contracts, error) {
	params, err := storage.Params(tx)
	if err != nil {
		return nil, err
	}

	contractMap := make(map[string]string)
	for _, param := range params {
		contractMap[param.Name] = param.Value
	}

	return contracts.NewByProxyAddress(client, contractMap["PROXY"])
}

func ToStorage(tx *sql.Tx, contracts *contracts.Contracts) error {
	if err := (storage.Param{"PROXY", contracts.ProxyAddress}).InsertOrUpdate(tx); err != nil {
		return err
	}
	return nil
}
