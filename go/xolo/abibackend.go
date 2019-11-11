package xolo

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
)

type AbiBackend struct {
	signerUrl   string
	client      *ethclient.Client
	transactOps *bind.TransactOpts
	networkID   int64
	lastNonce   int64
	pool        string
	pollTime    time.Duration
}

func NewAbiBackend(rpcUrl, signerUrl, pool string) (*AbiBackend, error) {
	conn, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, errors.Wrap(err, "NewAbiBackend-dial")
	}

	xbackend := AbiBackend{
		signerUrl: signerUrl,
		client:    conn,
		pool:      pool,
		pollTime:  time.Second * 2,
	}

	xbackend.transactOps = &bind.TransactOpts{
		From: common.HexToAddress("0x0000000000000000000000000000000000000000"),
		Signer: func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			return tx, nil
		},
	}

	nonce, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return nil, errors.Wrap(err, "NewAbiBackend-RandInt")
	}

	xbackend.lastNonce = nonce.Int64()

	netID, err := conn.NetworkID(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "NewAbiBackend-netid")
	}
	xbackend.networkID = netID.Int64()

	return &xbackend, nil
}

func (s *AbiBackend) TranslateTx(ctx context.Context, txhash common.Hash) (*common.Hash, error) {

	url := s.signerUrl + "/tx/" + txhash.Hex()
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrapf(err, "querying tx %v", txhash.Hex())
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("errcode %v querying tx %v", resp.StatusCode, txhash.Hex())
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "HTTP cannot read body %v", string(body))
	}
	var result XoloApiTranslateTxResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, errors.Wrapf(err, "HTTP cannot read body %v", string(body))
	}
	if !result.Success {
		return nil, fmt.Errorf("HTTP failed %v", result.Error)
	}
	if result.Tx != nil {
		tx := common.HexToHash(*result.Tx)
		return &tx, nil
	}
	return nil, nil
}

func (s *AbiBackend) WaitReceipt(ctx context.Context, tx *types.Transaction) (*types.Receipt, error) {
	var err error
	var txhash *common.Hash
	for txhash == nil {
		txhash, err = s.TranslateTx(ctx, tx.Hash())
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to translate tx %v", txhash.Hex())
		}
		if txhash == nil {
			select {
			case <-time.After(s.pollTime):
			case <-ctx.Done():
				return nil, fmt.Errorf("Context cancelled %v", txhash.Hex())
			}
		}
	}
	var receipt *types.Receipt
	for receipt == nil {
		receipt, err = s.client.TransactionReceipt(ctx, *txhash)
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

func (s *AbiBackend) serverSendTx(ctx context.Context, tx *types.Transaction) error {

	txr := xoloTx{
		Pool:    s.pool,
		TxID:    tx.Hash().Hex(),
		ChainID: s.networkID,
		To:      (*tx.To()).Hex(),
		Data:    hex.EncodeToString(tx.Data()),
		Value:   "0x" + tx.Value().Text(16),
	}
	txrjson, err := json.Marshal(&txr)
	if err != nil {
		return errors.Wrap(err, "serverSendTx-Marshal")
	}
	url := s.signerUrl + "/tx"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(txrjson))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "serverSendTx-do")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP bad code %v", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrapf(err, "HTTP cannot read body %v", string(body))
	}

	var result XoloApiBaseResult
	if err := json.Unmarshal(body, &result); err != nil {
		return errors.Wrapf(err, "HTTP cannot read body %v", string(body))
	}
	if !result.Success {
		return fmt.Errorf("HTTP failed %v", result.Error)
	}

	return nil
}

func (s *AbiBackend) TransactOps() *bind.TransactOpts {
	return s.transactOps
}
func (s *AbiBackend) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	return s.client.CodeAt(ctx, contract, blockNumber)
}
func (s *AbiBackend) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	return s.client.CallContract(ctx, call, blockNumber)
}
func (s *AbiBackend) PendingCodeAt(ctx context.Context, contract common.Address) ([]byte, error) {
	return s.client.PendingCodeAt(ctx, contract)
}
func (s *AbiBackend) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	return uint64(s.lastNonce), nil
}
func (s *AbiBackend) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return big.NewInt(1000000000), nil
}
func (s *AbiBackend) EstimateGas(ctx context.Context, call ethereum.CallMsg) (gas uint64, err error) {
	return 200000, nil
}
func (s *AbiBackend) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	s.lastNonce++
	return s.serverSendTx(ctx, tx)
}
func (s *AbiBackend) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	return s.client.FilterLogs(ctx, query)

}
func (s *AbiBackend) SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	return s.client.SubscribeFilterLogs(ctx, query, ch)
}
