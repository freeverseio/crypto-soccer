package xolo

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthClientFactory interface {
	Dial(URL string) (EthClient, error)
}

type EthClient interface {
	CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error)
	CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error)
	PendingCodeAt(ctx context.Context, contract common.Address) ([]byte, error)
	SendTransaction(ctx context.Context, tx *types.Transaction) error
	NetworkID(ctx context.Context) (*big.Int, error)
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
	EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error)
}

type GethEthClientFactory struct{}

func NewGethClientFactory() GethEthClientFactory {
	return GethEthClientFactory{}
}

func (g GethEthClientFactory) Dial(URL string) (EthClient, error) {
	return ethclient.Dial(URL)
}
