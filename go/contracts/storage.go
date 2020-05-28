package contracts

import (
	"database/sql"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/freeverseio/crypto-soccer/go/storage"
)

func NewFromStorage(client *ethclient.Client, tx *sql.Tx) (*Contracts, error) {
	params, err := storage.Params(tx)
	if err != nil {
		return nil, err
	}

	contractMap := make(map[string]string)
	for _, param := range params {
		contractMap[param.Name] = param.Value
	}

	return NewByProxyAddress(client, contractMap[ProxyName])
}

func (b Contracts) ToStorage(tx *sql.Tx) error {
	if err := (storage.Param{ProxyName, b.ProxyAddress}).InsertOrUpdate(tx); err != nil {
		return err
	}
	return nil
}
