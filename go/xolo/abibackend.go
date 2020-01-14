package xolo

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/prometheus/common/log"

	"github.com/pkg/errors"
)

var (
	agnosticAddr = common.HexToAddress("0xCAFECAFECAFECAFECAFECAFECAFECAFECAFECAFE")
)

type Signer struct {
	URL string
}

func NewSigner(URL string) *Signer {
	return &Signer{
		URL: URL,
	}
}

func (s *Signer) SendTransaction(ctx context.Context, tx *types.Transaction) (*common.Hash, error) {
	txr := xoloTx{
		To:    (*tx.To()).Hex(),
		Data:  hex.EncodeToString(tx.Data()),
		Value: "0x" + tx.Value().Text(16),
	}
	txrjson, err := json.Marshal(&txr)
	if err != nil {
		return nil, errors.Wrap(err, "serverSendTx-Marshal")
	}
	url := s.URL + "/tx"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(txrjson))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "serverSendTx-do")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP bad code %v", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "HTTP cannot read body %v", string(body))
	}

	var result XoloSendTxResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, errors.Wrapf(err, "HTTP cannot read body %v", string(body))
	}

	if !result.Success {
		return nil, fmt.Errorf("Call failed %v", result.Error)
	}

	txhash := common.HexToHash(*result.TxHash)
	return &txhash, nil
}

type SignerHA struct {
	signers   []*Signer
	available []time.Time
}

func NewSignerHA(urls []string) *SignerHA {
	available := []time.Time{}
	signers := []*Signer{}

	for _, url := range urls {
		signers = append(signers, NewSigner(url))
		available = append(available, time.Unix(0, 0))
	}
	return &SignerHA{
		signers:   signers,
		available: available,
	}
}

func (s *SignerHA) SendTransaction(ctx context.Context, tx *types.Transaction) (*common.Hash, error) {
	n := len(s.signers)
	i := rand.Intn(n)

	for n > 0 {
		if s.available[i].Before(time.Now()) {
			hash, err := s.signers[i].SendTransaction(ctx, tx)
			if err == nil {
				return hash, nil
			}
			log.Errorf("Unable to send transaction to %s : %s", s.signers[i].URL, err)
			s.available[i] = s.available[i].Add(time.Minute)
		}
		i = (i + 1) % len(s.signers)
		n--
	}
	return nil, errors.New("unable to process the transactions in any signer")
}

type ClientHA struct {
	sync.Mutex
	urls      []string
	clients   []*ethclient.Client
	available []time.Time
}

func NewClientHA(urls []string) *ClientHA {
	available := []time.Time{}
	clients := []*ethclient.Client{}

	for range urls {
		available = append(available, time.Unix(0, 0))
		clients = append(clients, nil)
	}

	return &ClientHA{
		urls:      urls,
		clients:   clients,
		available: available,
	}
}

func (s *ClientHA) doHA(opname string, f func(*ethclient.Client) (interface{}, error)) (interface{}, error) {
	n := len(s.clients)
	i := rand.Intn(n)

	for n > 0 {
		if s.available[i].Before(time.Now()) {
			var err error

			s.Lock()
			if s.clients[i] == nil {
				log.Info("Connecting to ", s.urls[i])
				s.clients[i], err = ethclient.Dial(s.urls[i])
			}
			s.Unlock()

			if err == nil {
				var res interface{}
				if res, err = f(s.clients[i]); err == nil || err == ethereum.NotFound {
					return res, nil
				}
			}

			log.Errorf("Unable to exec %s to %s : %s", opname, s.urls[i], err)
			s.clients[i] = nil
			s.available[i] = s.available[i].Add(time.Minute)
		}

		i = (i + 1) % len(s.clients)
		n--
	}
	return nil, errors.New("unable to process CallContract in any client")
}

func (s *ClientHA) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {

	res, err := s.doHA("CallContract", func(c *ethclient.Client) (interface{}, error) {
		return c.CallContract(ctx, call, blockNumber)
	})

	if err != nil {
		return nil, err
	}
	return res.([]byte), nil
}

func (s *ClientHA) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	res, err := s.doHA("TransactionReceipt", func(c *ethclient.Client) (interface{}, error) {
		return c.TransactionReceipt(ctx, txHash)
	})

	if err != nil {
		return nil, err
	}
	return res.(*types.Receipt), err
}

func (s *ClientHA) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	res, err := s.doHA("FilterLogs", func(c *ethclient.Client) (interface{}, error) {
		return c.FilterLogs(ctx, query)
	})

	if err != nil {
		return nil, err
	}
	return res.([]types.Log), err
}

func (s *ClientHA) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	res, err := s.doHA("CodeAt", func(c *ethclient.Client) (interface{}, error) {
		return c.CodeAt(ctx, contract, blockNumber)
	})

	if err != nil {
		return nil, err
	}
	return res.([]byte), err
}
func (s *ClientHA) PendingCodeAt(ctx context.Context, contract common.Address) ([]byte, error) {
	res, err := s.doHA("PendingCodeAt", func(c *ethclient.Client) (interface{}, error) {
		return c.PendingCodeAt(ctx, contract)
	})

	if err != nil {
		return nil, err
	}
	return res.([]byte), err
}

type AbiBackend struct {
	signer *SignerHA
	client *ClientHA

	transactOps *bind.TransactOpts
	lastNonce   int64

	TxMap    *TxMap
	pollTime time.Duration
}

func NewAbiBackend(rpcUrls, signerUrls []string) (*AbiBackend, error) {
	xbackend := AbiBackend{
		signer:   NewSignerHA(signerUrls),
		client:   NewClientHA(rpcUrls),
		pollTime: time.Second * 1,
		TxMap:    NewTxMap(),
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

func (s *AbiBackend) WaitReceipt(ctx context.Context, txint *types.Transaction) (*types.Receipt, error) {
	txhash := s.TxMap.Lookup(txint.Hash())
	if txhash == nil {
		return nil, fmt.Errorf("Cannot map tx")
	}

	var receipt *types.Receipt
	var err error

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

func (s *AbiBackend) TransactOps() *bind.TransactOpts {
	return s.transactOps
}

func (s *AbiBackend) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	return s.client.CallContract(ctx, call, blockNumber)
}

func (s *AbiBackend) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	return s.client.FilterLogs(ctx, query)
}

func (s *AbiBackend) SendTransaction(ctx context.Context, tx *types.Transaction) error {

	s.lastNonce++

	txhash, err := s.signer.SendTransaction(ctx, tx)
	if err != nil {
		return err
	}
	s.TxMap.Add(tx.Hash(), *txhash)

	return nil
}

func (s *AbiBackend) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	return s.client.CodeAt(ctx, contract, blockNumber)
}
func (s *AbiBackend) PendingCodeAt(ctx context.Context, contract common.Address) ([]byte, error) {
	return s.client.PendingCodeAt(ctx, contract)
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
