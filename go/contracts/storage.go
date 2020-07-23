package contracts

import (
	"database/sql"
	"errors"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/freeverseio/crypto-soccer/go/storage"
)

func NewFromStorage(client *ethclient.Client, tx *sql.Tx) (*Contracts, error) {
	leaguesParam, err := storage.ParamByName(tx, LeaguesName)
	if err != nil {
		return nil, err
	}
	if leaguesParam == nil {
		return nil, errors.New("no league contract in storage")
	}
	assetsParam, err := storage.ParamByName(tx, AssetsName)
	if err != nil {
		return nil, err
	}
	if assetsParam == nil {
		return nil, errors.New("no assets contract in storage")
	}
	evolutionParam, err := storage.ParamByName(tx, EvolutionName)
	if err != nil {
		return nil, err
	}
	if evolutionParam == nil {
		return nil, errors.New("no evolution contract in storage")
	}
	engineParam, err := storage.ParamByName(tx, EngineName)
	if err != nil {
		return nil, err
	}
	if engineParam == nil {
		return nil, errors.New("no engine contract in storage")
	}
	engineprecompParam, err := storage.ParamByName(tx, EnginePreCompName)
	if err != nil {
		return nil, err
	}
	if engineprecompParam == nil {
		return nil, errors.New("no engineprecomp contract in storage")
	}
	updatesParam, err := storage.ParamByName(tx, UpdatesName)
	if err != nil {
		return nil, err
	}
	if updatesParam == nil {
		return nil, errors.New("no updates contract in storage")
	}
	marketParam, err := storage.ParamByName(tx, MarketName)
	if err != nil {
		return nil, err
	}
	if marketParam == nil {
		return nil, errors.New("no market contract in storage")
	}
	utilsParam, err := storage.ParamByName(tx, UtilsName)
	if err != nil {
		return nil, err
	}
	if utilsParam == nil {
		return nil, errors.New("no utils contract in storage")
	}
	playandevolveParam, err := storage.ParamByName(tx, PlayAndEvolveName)
	if err != nil {
		return nil, err
	}
	if playandevolveParam == nil {
		return nil, errors.New("no playandevolve contract in storage")
	}
	shopParam, err := storage.ParamByName(tx, ShopName)
	if err != nil {
		return nil, err
	}
	if shopParam == nil {
		return nil, errors.New("no shop contract in storage")
	}
	trainingpointsParam, err := storage.ParamByName(tx, TrainingPointsName)
	if err != nil {
		return nil, err
	}
	if trainingpointsParam == nil {
		return nil, errors.New("no trainingpoints contract in storage")
	}
	constantsgettersParam, err := storage.ParamByName(tx, ConstantsGettersName)
	if err != nil {
		return nil, err
	}
	if constantsgettersParam == nil {
		return nil, errors.New("no constantsgetters contract in storage")
	}
	privilegedParam, err := storage.ParamByName(tx, PrivilegedName)
	if err != nil {
		return nil, err
	}
	if privilegedParam == nil {
		return nil, errors.New("no privileged contract in storage")
	}
	stakersParam, err := storage.ParamByName(tx, StakersName)
	if err != nil {
		return nil, err
	}
	if stakersParam == nil {
		return nil, errors.New("no stakers contract in storage")
	}
	directoryParam, err := storage.ParamByName(tx, DirectoryName)
	if err != nil {
		return nil, err
	}
	if directoryParam == nil {
		return nil, errors.New("no directory contract in storage")
	}
	proxyParam, err := storage.ParamByName(tx, ProxyName)
	if err != nil {
		return nil, err
	}
	if proxyParam == nil {
		return nil, errors.New("no proxy contract in storage")
	}

	return New(
		client,
		leaguesParam.Value,
		assetsParam.Value,
		evolutionParam.Value,
		engineParam.Value,
		engineprecompParam.Value,
		updatesParam.Value,
		marketParam.Value,
		utilsParam.Value,
		playandevolveParam.Value,
		shopParam.Value,
		trainingpointsParam.Value,
		constantsgettersParam.Value,
		privilegedParam.Value,
		stakersParam.Value,
		directoryParam.Value,
		proxyParam.Value,
	)
}

func (b Contracts) ToStorage(tx *sql.Tx) error {
	if err := (storage.Param{LeaguesName, b.LeaguesAddress}).InsertOrUpdate(tx); err != nil {
		return err
	}
	if err := (storage.Param{AssetsName, b.AssetsAddress}).InsertOrUpdate(tx); err != nil {
		return err
	}
	if err := (storage.Param{EvolutionName, b.EvolutionAddress}).InsertOrUpdate(tx); err != nil {
		return err
	}
	if err := (storage.Param{EngineName, b.EngineAddress}).InsertOrUpdate(tx); err != nil {
		return err
	}
	if err := (storage.Param{EnginePreCompName, b.EngineprecompAddress}).InsertOrUpdate(tx); err != nil {
		return err
	}
	if err := (storage.Param{UpdatesName, b.UpdatesAddress}).InsertOrUpdate(tx); err != nil {
		return err
	}
	if err := (storage.Param{MarketName, b.MarketAddress}).InsertOrUpdate(tx); err != nil {
		return err
	}
	if err := (storage.Param{UtilsName, b.UtilsAddress}).InsertOrUpdate(tx); err != nil {
		return err
	}
	if err := (storage.Param{PlayAndEvolveName, b.PlayandevolveAddress}).InsertOrUpdate(tx); err != nil {
		return err
	}
	if err := (storage.Param{ShopName, b.ShopAddress}).InsertOrUpdate(tx); err != nil {
		return err
	}
	if err := (storage.Param{TrainingPointsName, b.TrainingpointsAddress}).InsertOrUpdate(tx); err != nil {
		return err
	}
	if err := (storage.Param{ConstantsGettersName, b.ConstantsgettersAddress}).InsertOrUpdate(tx); err != nil {
		return err
	}
	if err := (storage.Param{PrivilegedName, b.PrivilegedAddress}).InsertOrUpdate(tx); err != nil {
		return err
	}
	if err := (storage.Param{StakersName, b.StakersAddress}).InsertOrUpdate(tx); err != nil {
		return err
	}
	if err := (storage.Param{DirectoryName, b.DirectoryAddress}).InsertOrUpdate(tx); err != nil {
		return err
	}
	if err := (storage.Param{ProxyName, b.ProxyAddress}).InsertOrUpdate(tx); err != nil {
		return err
	}
	return nil
}
