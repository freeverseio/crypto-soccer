package helper

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
)

func toCallArg(msg ethereum.CallMsg) interface{} {
	arg := map[string]interface{}{
		"from": msg.From,
		"to":   msg.To,
	}
	if len(msg.Data) > 0 {
		arg["data"] = hexutil.Bytes(msg.Data)
	}
	if msg.Value != nil {
		arg["value"] = (*hexutil.Big)(msg.Value)
	}
	if msg.Gas != 0 {
		arg["gas"] = hexutil.Uint64(msg.Gas)
	}
	if msg.GasPrice != nil {
		arg["gasPrice"] = (*hexutil.Big)(msg.GasPrice)
	}
	return arg
}

type ParityBackend struct {
	rpc    *rpc.Client
	Client *ethclient.Client
}

func NewParityBackend(rpcUrl string) (*ParityBackend, error) {
	c, err := rpc.DialContext(context.Background(), rpcUrl)
	if err != nil {
		return nil, err
	}
	return &ParityBackend{
		rpc:    c,
		Client: ethclient.NewClient(c),
	}, nil
}

func (s *ParityBackend) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	return s.Client.CodeAt(ctx, contract, blockNumber)
}
func (s *ParityBackend) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	return s.Client.CallContract(ctx, call, blockNumber)
}
func (s *ParityBackend) PendingCodeAt(ctx context.Context, contract common.Address) ([]byte, error) {
	return s.Client.PendingCodeAt(ctx, contract)
}
func (s *ParityBackend) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	return s.Client.PendingNonceAt(ctx, account)
}
func (s *ParityBackend) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return s.Client.SuggestGasPrice(ctx)
}
func (s *ParityBackend) EstimateGas(ctx context.Context, call ethereum.CallMsg) (gas uint64, err error) {
	gas, err = s.Client.EstimateGas(ctx, call)
	if err != nil && strings.Contains(fmt.Sprintf("%#v", err), `Data:"Reverted"`) {
		var hex map[string]interface{}
		err := s.rpc.CallContext(ctx, &hex, "trace_call", toCallArg(call), []string{"trace"})
		if err != nil {
			return 0, err
		}
		encoded := common.Hex2Bytes(hex["output"].(string)[2:])
		revertReason := string(encoded[4+32+32:])
		return 0, fmt.Errorf("Reverted with reason: %s", revertReason)
	}
	return gas, err
}
func (s *ParityBackend) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	return s.Client.SendTransaction(ctx, tx)
}
func (s *ParityBackend) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	return s.Client.FilterLogs(ctx, query)
}
func (s *ParityBackend) SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	return s.Client.SubscribeFilterLogs(ctx, query, ch)
}

func (s *ParityBackend) Transactor(keyAddr common.Address) *bind.TransactOpts {
	return &bind.TransactOpts{
		From: keyAddr,
		Signer: func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			if address != keyAddr {
				return nil, errors.New("not authorized to sign this account")
			}
			call := ethereum.CallMsg{
				From:     address,
				To:       tx.To(),
				Gas:      tx.Gas(),
				GasPrice: tx.GasPrice(),
				Value:    tx.Value(),
				Data:     tx.Data(),
			}

			var hex map[string]interface{}
			err := s.rpc.CallContext(context.Background(), &hex, "eth_signTransaction", toCallArg(call))
			if err != nil {
				return nil, err
			}
			rlptx, err := hexutil.Decode(hex["raw"].(string))
			if err != nil {
				return nil, err
			}
			var signedTx types.Transaction
			if err != signedTx.DecodeRLP(rlp.NewStream(bytes.NewReader(rlptx), math.MaxUint64)) {
				return nil, err
			}

			return &signedTx, nil
		},
	}
}
