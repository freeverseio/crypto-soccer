package xolo

import (
	"context"
	"encoding/hex"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type Server struct {
	sync.Mutex

	networkId *big.Int
	signer    *bind.TransactOpts
	eth       *ethclient.Client
}

func NewServer(signer *bind.TransactOpts, eth *ethclient.Client) (*Server, error) {

	networkId, err := eth.NetworkID(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "retriving network id")
	}

	return &Server{
		signer:    signer,
		eth:       eth,
		networkId: networkId,
	}, nil
}

func (x *Server) HttpPostTx(c *gin.Context) {
	var xtx xoloTx
	if err := c.ShouldBindJSON(&xtx); err != nil {
		c.JSON(200, XoloSendTxResult{Success: false, Error: "cannot parse htt body"})
		return
	}

	txhash, err := x.sendTx(&xtx)
	if err == nil {
		txhashstr := txhash.String()
		c.JSON(200, XoloSendTxResult{Success: true, TxHash: &txhashstr})
	} else {
		c.JSON(200, XoloSendTxResult{Success: false, Error: err.Error()})
	}
}

func (x *Server) sendTx(xtx *xoloTx) (*common.Hash, error) {

	ctx := context.Background()

	gasPrice, err := x.eth.SuggestGasPrice(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "sendTx-SuggestGasPrice")
	}
	data, err := hex.DecodeString(xtx.Data)
	if err != nil {
		return nil, errors.Wrap(err, "sendTx-DecodeString")
	}
	value, err := hexutil.DecodeBig(xtx.Value)
	if err != nil {
		return nil, errors.Wrap(err, "sendTx-DecodeBig")
	}
	to := common.HexToAddress(xtx.To)
	msg := ethereum.CallMsg{
		From:     x.signer.From,
		To:       &to,
		Gas:      8000000,
		GasPrice: gasPrice,
		Value:    value,
		Data:     data,
	}

	nonce, err := x.eth.PendingNonceAt(ctx, x.signer.From)
	if err != nil {
		return nil, errors.Wrap(err, "sendTx-PendingNonceAt")
	}

	msg.Gas, err = x.eth.EstimateGas(ctx, msg)
	if err != nil {
		return nil, errors.Wrap(err, "sendTx-EstimateGas")
	}

	tx := types.NewTransaction(nonce, *msg.To, msg.Value, msg.Gas, msg.GasPrice, msg.Data)
	ethsigner := types.NewEIP155Signer(x.networkId)
	x.Lock()
	tx, err = x.signer.Signer(ethsigner, x.signer.From, tx)
	x.Unlock()
	if err != nil {
		return nil, errors.Wrap(err, "sendTx-Signer")
	}

	if err := x.eth.SendTransaction(ctx, tx); err != nil {
		return nil, errors.Wrap(err, "sendTx-SendTransaction")
	}

	txhash := tx.Hash()

	return &txhash, nil
}
