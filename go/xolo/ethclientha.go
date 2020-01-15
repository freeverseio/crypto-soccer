package xolo

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/prometheus/common/log"
)

type EthClientHA struct {
	sync.Mutex
	factory   EthClientFactory
	urls      []string
	clients   []EthClient
	available []time.Time
	Rand      RandNFunc
}

func NewEthClientHA(factory EthClientFactory, urls []string) *EthClientHA {
	available := []time.Time{}
	clients := []EthClient{}

	for range urls {
		available = append(available, time.Unix(0, 0))
		clients = append(clients, nil)
	}

	return &EthClientHA{
		factory:   factory,
		urls:      urls,
		clients:   clients,
		available: available,
		Rand:      rand.Intn,
	}
}

func (s *EthClientHA) doHA(opname string, f func(EthClient) (interface{}, error)) (interface{}, error) {
	n := len(s.clients)
	i := s.Rand(n)

	for n > 0 {
		if s.available[i].Before(time.Now()) {
			var err error

			fmt.Println("ethclient-i->>>", i)

			s.Lock()
			if s.clients[i] == nil {
				log.Info("Connecting to ", s.urls[i])
				s.clients[i], err = s.factory.Dial(s.urls[i])
			}
			s.Unlock()

			if err == nil {
				var res interface{}
				if res, err = f(s.clients[i]); err == nil || err == ethereum.NotFound {
					return res, nil
				}
			}

			log.Errorf("Unable to exec %s to %s : %s", opname, s.urls[i], err)
			s.Lock()
			s.clients[i] = nil
			s.Unlock()
			s.available[i] = s.available[i].Add(time.Minute)
		}

		i = (i + 1) % len(s.clients)
		n--
	}
	return nil, errors.New("unable to process CallContract in any client")
}

func (s *EthClientHA) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {

	res, err := s.doHA("CallContract", func(c EthClient) (interface{}, error) {
		return c.CallContract(ctx, call, blockNumber)
	})

	if err != nil {
		return nil, err
	}
	return res.([]byte), nil
}

func (s *EthClientHA) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	res, err := s.doHA("TransactionReceipt", func(c EthClient) (interface{}, error) {
		return c.TransactionReceipt(ctx, txHash)
	})

	if err != nil {
		return nil, err
	}
	return res.(*types.Receipt), err
}

func (s *EthClientHA) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	res, err := s.doHA("FilterLogs", func(c EthClient) (interface{}, error) {
		return c.FilterLogs(ctx, query)
	})

	if err != nil {
		return nil, err
	}
	return res.([]types.Log), err
}

func (s *EthClientHA) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	res, err := s.doHA("CodeAt", func(c EthClient) (interface{}, error) {
		return c.CodeAt(ctx, contract, blockNumber)
	})

	if err != nil {
		return nil, err
	}
	return res.([]byte), err
}
func (s *EthClientHA) PendingCodeAt(ctx context.Context, contract common.Address) ([]byte, error) {
	res, err := s.doHA("PendingCodeAt", func(c EthClient) (interface{}, error) {
		return c.PendingCodeAt(ctx, contract)
	})

	if err != nil {
		return nil, err
	}
	return res.([]byte), err
}
