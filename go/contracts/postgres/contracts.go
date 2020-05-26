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

	return contracts.New(
		client,
		contractMap["LEAGUES"],
		contractMap["ASSETS"],
		contractMap["EVOLUTION"],
		contractMap["ENGINE"],
		contractMap["ENGINEPRECOMP"],
		contractMap["UPDATES"],
		contractMap["MARKET"],
		contractMap["UTILS"],
		contractMap["PLAYANDEVOLVE"],
		contractMap["SHOP"],
		contractMap["TRAININGPOINTS"],
		contractMap["CONSTANTSGETTERS"],
		contractMap["PRIVILEGED"],
		contractMap["STAKERS"],
		contractMap["DIRECTORY"],
	)
}
