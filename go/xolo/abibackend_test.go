package xolo

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"net/http"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"

	"github.com/freeverseio/crypto-soccer/go/xolo/example/abigen"
)

var (
	selectorInc, _ = hex.DecodeString("371303c0")
	selectorI, _   = hex.DecodeString("e5aa3d58")
)

func BigIntToBytes(b *big.Int) []byte {
	hexs := fmt.Sprintf("%x", b)
	for len(hexs) < 64 {
		hexs = "0" + hexs
	}
	v, _ := hex.DecodeString(hexs)
	return v
}

func BytesToBigInt(b []byte) *big.Int {
	big := big.NewInt(0)
	big.SetBytes(b)
	return big
}

func RandAddress() common.Address {
	addressBytes := make([]byte, common.AddressLength)
	rand.Read(addressBytes)
	return common.BytesToAddress(addressBytes)
}

type SmartContract interface {
	Fallback(view bool, from common.Address, data []byte) ([]byte, error)
}

type CounterSmartContract struct {
	i *big.Int
}

func NewCounterSmartContract() *CounterSmartContract {
	return &CounterSmartContract{
		i: big.NewInt(0),
	}
}

func (c *CounterSmartContract) Fallback(view bool, from common.Address, data []byte) ([]byte, error) {
	selector := data[:4]

	if bytes.Compare(selector, selectorInc) == 0 {
		return c.Inc(view, from, data)
	}
	if bytes.Compare(selector, selectorI) == 0 {
		return c.I(view, from, data)
	}
	return nil, fmt.Errorf("invalid selector %x", selector)
}

func (c *CounterSmartContract) I(view bool, from common.Address, data []byte) ([]byte, error) {
	return BigIntToBytes(c.i), nil
}

func (c *CounterSmartContract) Inc(view bool, from common.Address, data []byte) ([]byte, error) {
	c.i = big.NewInt(c.i.Int64() + 1)
	return []byte{}, nil
}

type Blockchain struct {
	nextNonce      map[common.Address]uint64
	smartcontracts map[common.Address]SmartContract
	pool           []*types.Transaction
	rcpt           map[common.Hash]*types.Receipt
}

func NewBlockchain() *Blockchain {
	return &Blockchain{
		nextNonce:      map[common.Address]uint64{},
		smartcontracts: map[common.Address]SmartContract{},
		pool:           []*types.Transaction{},
		rcpt:           map[common.Hash]*types.Receipt{},
	}
}

func (b *Blockchain) Deploy(sc SmartContract) common.Address {
	address := RandAddress()
	b.smartcontracts[address] = sc
	return address
}

func (b *Blockchain) incNonce(address common.Address) {
	if _, ok := b.nextNonce[address]; ok {
		b.nextNonce[address]++
	} else {
		b.nextNonce[address] = 1
	}
}

func (b *Blockchain) NextNonce(address common.Address) uint64 {
	if nonce, ok := b.nextNonce[address]; ok {
		return nonce
	}
	return 0
}

func (b *Blockchain) AddTx(tx *types.Transaction) {
	b.pool = append(b.pool, tx)
}

func (b *Blockchain) Call(call ethereum.CallMsg) ([]byte, error) {
	if sc, ok := b.smartcontracts[*call.To]; ok {
		return sc.Fallback(true, call.From, call.Data)
	}
	return nil, errors.New("Invalid SC address")
}

func (b *Blockchain) execTx(tx *types.Transaction) (*types.Receipt, error) {
	from := b.extractFrom(tx)
	if sc, ok := b.smartcontracts[*tx.To()]; ok {
		_, err := sc.Fallback(false, from, tx.Data())
		if err != nil {
			return nil, err
		}
		b.incNonce(from)
		var bloom [types.BloomByteLength]byte
		rcpt := types.Receipt{
			PostState:         nil,
			Status:            types.ReceiptStatusSuccessful,
			CumulativeGasUsed: 1000,
			Bloom:             types.BytesToBloom(bloom[:]),
			Logs:              nil,
			TxHash:            tx.Hash(),
			ContractAddress:   *tx.To(),
			GasUsed:           1000,
			BlockHash:         tx.Hash(),
			BlockNumber:       big.NewInt(1),
			TransactionIndex:  1,
		}
		return &rcpt, err
	}
	return nil, errors.New("Invalid SC address")
}

func (b *Blockchain) extractFrom(tx *types.Transaction) common.Address {
	msg, err := tx.AsMessage(types.NewEIP155Signer(tx.ChainId()))
	if err != nil {
		log.Fatal(err)
	}
	return msg.From()
}

func (b *Blockchain) TransactionReceipt(txHash common.Hash) *types.Receipt {
	rcpt, _ := b.rcpt[txHash]
	return rcpt
}

func (b *Blockchain) Mine() error {
	var err error
	anyProcessed := true
	for anyProcessed {
		anyProcessed = false
		for n, tx := range b.pool {

			// check nonce
			from := b.extractFrom(tx)
			if b.NextNonce(from) != tx.Nonce() {
				continue
			}
			// run tx
			b.pool = append(b.pool[:n], b.pool[n+1:]...)
			if b.rcpt[tx.Hash()], err = b.execTx(tx); err != nil {
				return err
			}

			anyProcessed = true
			break
		}
	}
	return nil
}

type EthClientMock struct {
	URL        string
	blockchain *Blockchain
	factory    *EthClientMockFactory
}

var errNotAllowed = errors.New("Disabled")

func (m *EthClientMock) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	if !m.factory.allowed[m.URL] {
		return nil, errNotAllowed
	}
	m.factory.calls[m.URL] = append(m.factory.calls[m.URL], "CC")
	return m.blockchain.Call(call)
}
func (m *EthClientMock) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	if !m.factory.allowed[m.URL] {
		return nil, errNotAllowed
	}
	m.factory.calls[m.URL] = append(m.factory.calls[m.URL], "TR")
	if rcpt := m.blockchain.TransactionReceipt(txHash); rcpt != nil {
		return rcpt, nil
	}
	return nil, ethereum.NotFound
}
func (m *EthClientMock) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	if !m.factory.allowed[m.URL] {
		return nil, errNotAllowed
	}
	m.factory.calls[m.URL] = append(m.factory.calls[m.URL], "FL")
	return nil, nil
}
func (m *EthClientMock) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	if !m.factory.allowed[m.URL] {
		return nil, errNotAllowed
	}
	m.factory.calls[m.URL] = append(m.factory.calls[m.URL], "CA")
	return []byte{0}, nil
}
func (m *EthClientMock) PendingCodeAt(ctx context.Context, contract common.Address) ([]byte, error) {
	if !m.factory.allowed[m.URL] {
		return nil, errNotAllowed
	}
	m.factory.calls[m.URL] = append(m.factory.calls[m.URL], "PC")
	return []byte{0}, nil
}
func (m *EthClientMock) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	if !m.factory.allowed[m.URL] {
		return errNotAllowed
	}
	m.factory.calls[m.URL] = append(m.factory.calls[m.URL], "ST")
	m.blockchain.AddTx(tx)
	return m.blockchain.Mine()
}
func (m *EthClientMock) NetworkID(ctx context.Context) (*big.Int, error) {
	if !m.factory.allowed[m.URL] {
		return nil, errNotAllowed
	}
	m.factory.calls[m.URL] = append(m.factory.calls[m.URL], "NI")
	return big.NewInt(1), nil
}
func (m *EthClientMock) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	if !m.factory.allowed[m.URL] {
		return nil, errNotAllowed
	}
	m.factory.calls[m.URL] = append(m.factory.calls[m.URL], "GP")
	return big.NewInt(20000), nil
}
func (m *EthClientMock) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	if !m.factory.allowed[m.URL] {
		return 0, errNotAllowed
	}
	m.factory.calls[m.URL] = append(m.factory.calls[m.URL], "PN")
	return m.blockchain.NextNonce(account), nil
}
func (m *EthClientMock) EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
	if !m.factory.allowed[m.URL] {
		return 0, errNotAllowed
	}
	m.factory.calls[m.URL] = append(m.factory.calls[m.URL], "EG")
	return 20000, nil
}

func NewEthClientMockFactory(blockchain *Blockchain) *EthClientMockFactory {
	return &EthClientMockFactory{
		blockchain: blockchain,
		allowed:    map[string]bool{},
		calls:      map[string][]string{},
	}
}

type EthClientMockFactory struct {
	allowed    map[string]bool
	calls      map[string][]string
	blockchain *Blockchain
}

func (g *EthClientMockFactory) Dial(URL string) (EthClient, error) {
	if _, ok := g.allowed[URL]; !ok {
		g.allowed[URL] = true
		g.calls[URL] = []string{}
	}
	return &EthClientMock{
		URL:        URL,
		blockchain: g.blockchain,
		factory:    g,
	}, nil
}

const acc1key = `{"address":"be3a732e058fdfdb3457ba1bb1d87f9c200982f2","crypto":{"cipher":"aes-128-ctr","ciphertext":"8fe134c7059aebde9043f6454a8b6451d52b3d4e4c9162728fd35f1fee05c229","cipherparams":{"iv":"1eb6f75e10f5de357f9019176fd9a8d7"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"8f42b2a446065308d01778c16af2cea37f0829096de4796c0d81b1ed140817b5"},"mac":"227761a06aa2881ee51a7fcd7cb2f35317dbbc1ba9b4749e562c7eb045650ba2"},"id":"68c3a22d-aa36-49cb-a2e9-b97c717f3958","version":3}`
const acc1pwd = "11111111"
const acc2key = `{"address":"c210a199df85b19c89aa527a300597e3b0023be5","crypto":{"cipher":"aes-128-ctr","ciphertext":"ff73cd25bdd57e2b176bfe849175357181c816f5139c4166b2e6d96760e89afd","cipherparams":{"iv":"540c3a443300f2539f53af7fe5df81c1"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"d3658037ec4a799368cd1d77e50ff8347e4fb32cdc3379519d0006ae839a53ab"},"mac":"9c37d76c4f9c0bf887a3d0cb808e52d44120781e8503aa3703f97da76ac3297c"},"id":"b773cb46-ceee-4016-a87a-abad55439a63","version":3}`
const acc2pwd = "11111111"

func TestBlockchain(t *testing.T) {
	ethsigner := types.NewEIP155Signer(big.NewInt(1))
	acc1, err := bind.NewTransactor(strings.NewReader(acc1key), acc1pwd)

	blockchain := NewBlockchain()
	counter := blockchain.Deploy(NewCounterSmartContract())

	txinc1 := types.NewTransaction(blockchain.NextNonce(acc1.From), counter, big.NewInt(0), 0, big.NewInt(0), selectorInc)
	txinc1, err = acc1.Signer(ethsigner, acc1.From, txinc1)
	txinc2 := types.NewTransaction(blockchain.NextNonce(acc1.From)+1, counter, big.NewInt(0), 0, big.NewInt(0), selectorInc)
	txinc2, err = acc1.Signer(ethsigner, acc1.From, txinc2)

	blockchain.AddTx(txinc1)
	blockchain.AddTx(txinc2)
	assert.Nil(t, blockchain.Mine())

	two, err := blockchain.Call(ethereum.CallMsg{
		From:     acc1.From,
		To:       &counter,
		Gas:      0,
		GasPrice: big.NewInt(0),
		Value:    big.NewInt(0),
		Data:     selectorI,
	})
	assert.Nil(t, err)
	assert.Equal(t, int64(2), BytesToBigInt(two).Int64())
	assert.Equal(t, uint64(2), blockchain.NextNonce(acc1.From))
}

func TestBasic(t *testing.T) {

	client1Url := "1"

	blockchain := NewBlockchain()
	counterAddr := blockchain.Deploy(NewCounterSmartContract())
	factory := NewEthClientMockFactory(blockchain)

	client1, err := factory.Dial(client1Url)
	assert.Nil(t, err)
	client1signer, err := bind.NewTransactor(strings.NewReader(acc1key), acc1pwd)
	assert.Nil(t, err)

	xserver, err := NewServer(client1signer, client1)
	assert.Nil(t, err)

	engine := gin.Default()
	srv := &http.Server{
		Addr:    ":8004",
		Handler: engine,
	}
	engine.POST("/tx", xserver.HttpPostTx)
	go srv.ListenAndServe()

	xbackend, err := NewAbiBackend(
		factory,
		[]string{client1Url},
		[]string{"http://localhost:8004"})
	assert.Nil(t, err)

	counter, err := abigen.NewCounter(counterAddr, xbackend)
	assert.Nil(t, err)

	session := abigen.CounterSession{
		Contract:     counter,
		CallOpts:     bind.CallOpts{},
		TransactOpts: *xbackend.TransactOps(),
	}

	previousI, err := session.I()
	assert.Nil(t, err)

	size := 10

	txs := []*types.Transaction{}
	for n := 0; n < size; n++ {
		tx, err := session.Inc()
		assert.Nil(t, err)
		txs = append(txs, tx)
	}

	for n := 0; n < size; n++ {
		receipt, err := xbackend.WaitReceipt(context.Background(), txs[n])
		assert.Nil(t, err)
		if receipt.Status != types.ReceiptStatusSuccessful {
			t.Error("!ReceiptStatusSuccessful")
		}
	}

	nextI, err := session.I()
	assert.Nil(t, err)

	diff := nextI.Uint64() - previousI.Uint64()
	assert.Equal(t, diff, uint64(size))

	srv.Shutdown(context.Background())
}

type FakeRand struct {
	i int
}

func (r *FakeRand) IntN(n int) int {
	v := r.i % n
	r.i++
	return v
}

func TestHA(t *testing.T) {
	client1Url := "1"
	client2Url := "2"

	// blockchain
	blockchain := NewBlockchain()
	counterAddr := blockchain.Deploy(NewCounterSmartContract())
	factory := NewEthClientMockFactory(blockchain)

	// xolo servers
	client1, err := factory.Dial(client1Url)
	assert.Nil(t, err)
	client1signer, err := bind.NewTransactor(strings.NewReader(acc1key), acc1pwd)
	assert.Nil(t, err)

	client2, err := factory.Dial(client2Url)
	assert.Nil(t, err)
	client2signer, err := bind.NewTransactor(strings.NewReader(acc2key), acc2pwd)
	assert.Nil(t, err)

	xserver1, err := NewServer(client1signer, client1)
	assert.Nil(t, err)
	xserver2, err := NewServer(client2signer, client2)
	assert.Nil(t, err)

	engine1 := gin.Default()
	srv1 := &http.Server{
		Addr:    ":8004",
		Handler: engine1,
	}
	engine1.POST("/tx", xserver1.HttpPostTx)
	go srv1.ListenAndServe()

	engine2 := gin.Default()
	srv2 := &http.Server{
		Addr:    ":8005",
		Handler: engine2,
	}
	engine2.POST("/tx", xserver2.HttpPostTx)
	go srv2.ListenAndServe()

	// backend
	xbackend, err := NewAbiBackend(
		factory,
		[]string{client1Url, client2Url},
		[]string{"http://localhost:8004", "http://localhost:8005"})
	assert.Nil(t, err)

	var fakeRand FakeRand
	xbackend.Xolo().Rand = fakeRand.IntN
	xbackend.Eth().Rand = fakeRand.IntN

	counter, err := abigen.NewCounter(counterAddr, xbackend)
	assert.Nil(t, err)

	session := abigen.CounterSession{
		Contract:     counter,
		CallOpts:     bind.CallOpts{},
		TransactOpts: *xbackend.TransactOps(),
	}

	previousI, err := session.I()
	assert.Nil(t, err)

	size := 10

	// send transactions, stop first node at half
	txs := []*types.Transaction{}
	for n := 0; n < size/2; n++ {
		tx, err := session.Inc()
		assert.Nil(t, err)
		txs = append(txs, tx)
	}
	factory.allowed[client1Url] = false
	fact1calls, fact2calls := len(factory.calls[client1Url]), len(factory.calls[client2Url])
	for n := 0; n < size/2; n++ {
		tx, err := session.Inc()
		assert.Nil(t, err)
		txs = append(txs, tx)
	}
	assert.True(t, len(factory.calls[client1Url]) == fact1calls)
	assert.True(t, len(factory.calls[client2Url]) > fact2calls)

	// collect receipts, switch activated
	factory.allowed[client1Url] = true
	factory.allowed[client2Url] = false

	fact1calls, fact2calls = len(factory.calls[client1Url]), len(factory.calls[client2Url])

	for n := 0; n < size; n++ {
		receipt, err := xbackend.WaitReceipt(context.Background(), txs[n])
		assert.Nil(t, err)
		if receipt.Status != types.ReceiptStatusSuccessful {
			t.Error("!ReceiptStatusSuccessful")
		}
	}
	assert.True(t, len(factory.calls[client1Url]) > fact1calls)
	assert.True(t, len(factory.calls[client2Url]) == fact2calls)

	nextI, err := session.I()
	assert.Nil(t, err)

	diff := nextI.Uint64() - previousI.Uint64()
	assert.Equal(t, diff, uint64(size))

	// shutdown server
	srv1.Shutdown(context.Background())
	srv2.Shutdown(context.Background())
}
