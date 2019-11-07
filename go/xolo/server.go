package xolo

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type Signer struct {
	tx     *Tx
	signer *bind.TransactOpts
}

func (s *Signer) Address() common.Address {
	return s.signer.From
}

type RpcClient struct {
	failed int
	client *ethclient.Client
}

type SignerPool struct {
	queue       []*Tx
	signers     []*Signer
	rpcs        []*RpcClient
	waitReceipt time.Duration
}

type Tx struct {
	xtx  *xoloTx
	etx  *types.Transaction
	rcpt *types.Receipt
	err  error
	in   time.Time
	out  time.Time
}

type Server struct {
	sync.Mutex

	pools   map[string]*SignerPool
	pending map[common.Hash]*Tx
	done    map[common.Hash]*Tx

	waitOutQueue time.Duration
}

func NewServer() (*Server, error) {

	return &Server{
		pools:        make(map[string]*SignerPool),
		pending:      make(map[common.Hash]*Tx),
		done:         make(map[common.Hash]*Tx),
		waitOutQueue: time.Second * 20,
	}, nil
}

func (x *Server) AddPool(poolname string, waitReceipt time.Duration) error {

	if _, ok := x.pools[poolname]; ok {
		return errors.New("pool " + poolname + " already exist")
	}

	x.pools[poolname] = &SignerPool{
		signers:     []*Signer{},
		queue:       []*Tx{},
		waitReceipt: waitReceipt,
	}

	return nil
}

func (x *Server) AddRpcClient(poolname string, rpcclient *ethclient.Client) error {
	pool, ok := x.pools[poolname]
	if !ok {
		return errors.New("pool " + poolname + " does not exist")
	}
	pool.rpcs = append(pool.rpcs, &RpcClient{
		client: rpcclient,
	})
	return nil
}

func (x *Server) AddSigner(poolname string, key string, passwd string) (*Signer, error) {

	transactOps, err := bind.NewTransactor(strings.NewReader(key), passwd)
	if err != nil {
		return nil, errors.Wrap(err, "cannot add key")
	}

	pool, ok := x.pools[poolname]
	if !ok {
		return nil, errors.New("pool " + poolname + " does not exist")
	}

	signer := &Signer{
		signer: transactOps,
	}

	pool.signers = append(pool.signers, signer)
	x.pools[poolname] = pool

	return signer, nil
}

func (x *Server) Start(engine *gin.Engine) {
	go x.loop()
}

func (x *Server) ServeGetTx(c *gin.Context) {
	x.Lock()
	defer x.Unlock()

	txhash := common.HexToHash(c.Param("txhash"))
	if tx, ok := x.done[txhash]; ok {
		txhashhex := tx.etx.Hash().Hex()
		c.JSON(200, XoloApiTranslateTxResult{
			Success: true,
			Tx:      &txhashhex,
		})
	} else if _, ok = x.pending[txhash]; ok {
		c.JSON(200, XoloApiTranslateTxResult{
			Success: true,
			Tx:      nil,
		})
	} else {
		c.JSON(200, XoloApiBaseResult{false, "cannot find hash"})
	}

}

func (x *Server) ServePostTx(c *gin.Context) {
	var xtx xoloTx
	if err := c.ShouldBindJSON(&xtx); err != nil {
		c.JSON(200, XoloApiBaseResult{false, "cannot parse htt body"})
		return
	}

	x.Lock()
	pool, ok := x.pools[xtx.Pool]
	if ok {
		tx := &Tx{xtx: &xtx, in: time.Now()}
		pool.queue = append(pool.queue, tx)
		x.pending[common.HexToHash(xtx.TxID)] = tx
	}
	x.Unlock()

	if ok {
		c.JSON(200, XoloApiBaseResult{Success: true})
	} else {
		c.JSON(200, XoloApiBaseResult{Success: false, Error: "pool not found"})
	}
}

func (x *Server) Info() string {
	x.Lock()
	var buffer bytes.Buffer

	write := func(s string) {
		if buffer.Len() > 0 {
			buffer.WriteString("\n")
		}
		buffer.WriteString(s)
	}

	for poolname, p := range x.pools {
		for _, t := range p.queue {
			write(fmt.Sprintf("⏳ [%v] %v", poolname, t.xtx.TxID))
		}
		for i, s := range p.signers {
			write(fmt.Sprintf("🔨 [%v-%v] ", poolname, i))
			if s.tx == nil {
				buffer.WriteString(fmt.Sprintf(" <empty>"))
			} else {
				buffer.WriteString(fmt.Sprintf("%v", s.tx.xtx.TxID))
			}
		}
	}
	for _, t := range x.done {
		if t.err == nil {
			write(fmt.Sprintf("✅ %v", t.xtx.TxID))
		} else {
			write(fmt.Sprintf("❌ %v", t.xtx.TxID))
		}
	}
	x.Unlock()
	return buffer.String()
}

func (x *Server) loop() {
	for {
		x.Lock()

		// send txs
		for _, p := range x.pools {
			for _, s := range p.signers {
				if len(p.queue) > 0 && s.tx == nil {
					s.tx, p.queue = p.queue[0], p.queue[1:]
					go x.process(p, s)
				}
			}
		}

		// clean old
		for txh, tx := range x.done {
			if time.Now().Sub(tx.out) > x.waitOutQueue {
				delete(x.done, txh)
			}
		}
		x.Unlock()

		time.Sleep(time.Second)
	}
}

func (x *Server) randrpc(pool *SignerPool) *ethclient.Client {
	return pool.rpcs[rand.Intn(len(pool.rpcs))].client
}

func (x *Server) process(pool *SignerPool, signer *Signer) error {

	ctx := context.Background()
	xtx := signer.tx.xtx

	defer func() {
		x.Lock()
		signer.tx.out = time.Now()
		txhash := common.HexToHash(xtx.TxID)
		x.done[txhash] = signer.tx
		delete(x.pending, txhash)
		signer.tx = nil
		x.Unlock()
	}()

	onerr := func(err error, msg string) error {
		signer.tx.err = errors.Wrap(err, msg)
		return signer.tx.err
	}

	gasPrice, err := x.randrpc(pool).SuggestGasPrice(ctx)
	if err != nil {
		return onerr(err, "sendTx-SuggestGasPrice")
	}
	data, err := hex.DecodeString(xtx.Data)
	if err != nil {
		return onerr(err, "sendTx-DecodeString")
	}
	value, err := hexutil.DecodeBig(xtx.Value)
	if err != nil {
		return onerr(err, "sendTx-DecodeBig")
	}
	to := common.HexToAddress(xtx.To)
	msg := ethereum.CallMsg{
		From:     signer.signer.From,
		To:       &to,
		Gas:      8000000,
		GasPrice: gasPrice,
		Value:    value,
		Data:     data,
	}

	nonce, err := x.randrpc(pool).PendingNonceAt(ctx, signer.signer.From)
	if err != nil {
		return onerr(err, "sendTx-PendingNonceAt")
	}

	msg.Gas, err = x.randrpc(pool).EstimateGas(ctx, msg)
	if err != nil {
		return onerr(err, "sendTx-EstimateGas")
	}

	tx := types.NewTransaction(nonce, *msg.To, msg.Value, msg.Gas, msg.GasPrice, msg.Data)
	ethsigner := types.NewEIP155Signer(big.NewInt(xtx.ChainID))
	tx, err = signer.signer.Signer(ethsigner, signer.signer.From, tx)
	if err != nil {
		return onerr(err, "sendTx-Signer")
	}

	signer.tx.etx = tx
	if err := x.randrpc(pool).SendTransaction(ctx, tx); err != nil {
		return onerr(err, "sendTx-SendTransaction")
	}

	start := time.Now()
	for time.Now().Sub(start) < pool.waitReceipt {
		signer.tx.rcpt, err = x.randrpc(pool).TransactionReceipt(ctx, tx.Hash())
		if err == nil {
			break
		} else if err == ethereum.NotFound {
			time.Sleep(time.Second)
		} else {
			return onerr(err, "sendRx-Receipt")
		}
	}

	if signer.tx.rcpt == nil {
		return onerr(nil, "sendTx-rcpt-timeout")
	}

	return nil
}
