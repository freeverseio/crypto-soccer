package helper

import (
	"context"
	"errors"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func WaitReceipt(client *ethclient.Client, tx *types.Transaction, timeoutSec uint8) error {
	receiptTimeout := time.Second * time.Duration(timeoutSec)
	start := time.Now()
	ctx := context.TODO()
	var receipt *types.Receipt

	for receipt == nil && time.Now().Sub(start) < receiptTimeout {
		receipt, err := client.TransactionReceipt(ctx, tx.Hash())
		if err == nil && receipt != nil {
			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	return errors.New("Timeout waiting for receipt")
}

func WaitReceipts(client *ethclient.Client, txs []*types.Transaction, timeoutSec uint8) error {
	for _, tx := range txs {
		err := WaitReceipt(client, tx, timeoutSec)
		if err != nil {
			return err
		}
	}
	return nil
}
