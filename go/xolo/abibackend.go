package xolo

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"

	"github.com/pkg/errors"
)

var (
	agnosticAddr = common.HexToAddress("0xCAFECAFECAFECAFECAFECAFECAFECAFECAFECAFE")
)

type AbiBackend struct {
	xolo *XoloClientHA
	eth  *EthClientHA

	transactOps *bind.TransactOpts
	lastNonce   int64

	txMap    *TxMap
	pollTime time.Duration
}

func NewAbiBackend(factory EthClientFactory, rpcUrls, xoloUrls []string) (*AbiBackend, error) {
	xbackend := AbiBackend{
		xolo:     NewXoloClientHA(xoloUrls),
		eth:      NewEthClientHA(factory, rpcUrls),
		pollTime: time.Second * 1,
		txMap:    NewTxMap(),
	}

	xbackend.transactOps = &bind.TransactOpts{
		From: agnosticAddr,
		Signer: func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			return tx, nil
		},
	}

	xbackend.lastNonce = 0

	return &xbackend, nil
}

func (s *AbiBackend) Xolo() *XoloClientHA {
	return s.xolo
}
func (s *AbiBackend) Eth() *EthClientHA {
	return s.eth
}

func (s *AbiBackend) WaitReceipt(ctx context.Context, txint *types.Transaction) (*types.Receipt, error) {
	txhash := s.txMap.Lookup(txint.Hash())
	if txhash == nil {
		return nil, fmt.Errorf("Cannot map tx")
	}

	var receipt *types.Receipt
	var err error

	for receipt == nil {

		receipt, err = s.eth.TransactionReceipt(ctx, *txhash)
		if err == ethereum.NotFound {
			select {
			case <-time.After(s.pollTime):
			case <-ctx.Done():
				return nil, fmt.Errorf("Context cancelled %v", txhash.Hex())
			}
		} else if err != nil {
			return nil, errors.Wrapf(err, "Failed to query receipt tx %v", txhash.Hex())
		}

	}
	return receipt, nil
}

func (s *AbiBackend) TransactOps() *bind.TransactOpts {
	return s.transactOps
}

func (s *AbiBackend) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	return s.eth.CallContract(ctx, call, blockNumber)
}

func (s *AbiBackend) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	return s.eth.FilterLogs(ctx, query)
}

func (s *AbiBackend) SendTransaction(ctx context.Context, tx *types.Transaction) error {

	s.lastNonce++

	txhash, err := s.xolo.SendTransaction(ctx, tx)
	if err != nil {
		return err
	}
	s.txMap.Add(tx.Hash(), *txhash)

	return nil
}

func (s *AbiBackend) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	return s.eth.CodeAt(ctx, contract, blockNumber)
}
func (s *AbiBackend) PendingCodeAt(ctx context.Context, contract common.Address) ([]byte, error) {
	return s.eth.PendingCodeAt(ctx, contract)
}

func (s *AbiBackend) SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	return nil, errors.New("SubscribeFilterLogs not supported")
}

func (s *AbiBackend) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	if account != agnosticAddr {
		return 0, errors.New("Cannot PendingNonceAt with non-agnotic address")
	}
	return uint64(s.lastNonce), nil
}
func (s *AbiBackend) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return big.NewInt(params.GWei), nil
}
func (s *AbiBackend) EstimateGas(ctx context.Context, call ethereum.CallMsg) (gas uint64, err error) {
	return 200000, nil
}
