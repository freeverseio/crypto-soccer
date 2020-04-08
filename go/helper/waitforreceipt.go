package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func TransactionToString(tx *types.Transaction) (string, error) {
	b, err := json.MarshalIndent(tx, "", "\t")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func WaitReceipt(client *ethclient.Client, tx *types.Transaction, timeoutSec int) (*types.Receipt, error) {
	receiptTimeout := time.Second * time.Duration(timeoutSec)
	start := time.Now()
	ctx := context.TODO()
	var receipt *types.Receipt

	for receipt == nil && time.Now().Sub(start) < receiptTimeout {
		receipt, err := client.TransactionReceipt(ctx, tx.Hash())
		if err == nil && receipt != nil {
			return receipt, nil
		} else if err != nil {
			return nil, err
		}
		time.Sleep(200 * time.Millisecond)
	}

	dump, err := TransactionToString(tx)
	if err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("Timout in transaction:\n%s", dump)
}

func WaitReceipts(client *ethclient.Client, txs []*types.Transaction, timeoutSec int) error {
	for _, tx := range txs {
		_, err := WaitReceipt(client, tx, timeoutSec)
		if err != nil {
			return err
		}
	}
	return nil
}
