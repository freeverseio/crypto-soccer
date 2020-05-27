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
		contractMap["PROXY"],
	)
}

func ToStorage(tx *sql.Tx, contracts *contracts.Contracts) error {
	if err := (storage.Param{"LEAGUES", contracts.LeaguesAddress}).InsertOrUpdate(tx); err != nil {
		return err
	}
	if err := (storage.Param{"ASSETS", contracts.AssetsAddress}).InsertOrUpdate(tx); err != nil {
		return err
	}
	if err := (storage.Param{"EVOLUTION", contracts.EvolutionAddress}).InsertOrUpdate(tx); err != nil {
		return err
	}
	if err := (storage.Param{"ENGINE", contracts.EngineAddress}).InsertOrUpdate(tx); err != nil {
		return err
	}
	if err := (storage.Param{"ENGINEPRECOMP", contracts.EngineprecompAddress}).InsertOrUpdate(tx); err != nil {
		return err
	}
	if err := (storage.Param{"UPDATES", contracts.UpdatesAddress}).InsertOrUpdate(tx); err != nil {
		return err
	}
	if err := (storage.Param{"MARKET", contracts.MarketAddress}).InsertOrUpdate(tx); err != nil {
		return err
	}
	if err := (storage.Param{"UTILS", contracts.UtilsAddress}).InsertOrUpdate(tx); err != nil {
		return err
	}
	if err := (storage.Param{"PLAYANDEVOLVE", contracts.PlayandevolveAddress}).InsertOrUpdate(tx); err != nil {
		return err
	}
	if err := (storage.Param{"SHOP", contracts.ShopAddress}).InsertOrUpdate(tx); err != nil {
		return err
	}
	if err := (storage.Param{"TRAININGPOINTS", contracts.TrainingpointsAddress}).InsertOrUpdate(tx); err != nil {
		return err
	}
	if err := (storage.Param{"CONSTANTSGETTERS", contracts.ConstantsgettersAddress}).InsertOrUpdate(tx); err != nil {
		return err
	}
	if err := (storage.Param{"PRIVILEGED", contracts.PrivilegedAddress}).InsertOrUpdate(tx); err != nil {
		return err
	}
	if err := (storage.Param{"STAKERS", contracts.StakersAddress}).InsertOrUpdate(tx); err != nil {
		return err
	}
	if err := (storage.Param{"DIRECTORY", contracts.DirectoryAddress}).InsertOrUpdate(tx); err != nil {
		return err
	}
	return nil
}
