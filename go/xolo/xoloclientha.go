package xolo

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/prometheus/common/log"
)

type XoloClient struct {
	URL string
}

func NewXoloClient(URL string) *XoloClient {
	return &XoloClient{
		URL: URL,
	}
}

func (s *XoloClient) SendTransaction(ctx context.Context, tx *types.Transaction) (*common.Hash, error) {
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

type XoloClientHA struct {
	Rand RandNFunc

	XoloClients []*XoloClient
	available   []time.Time
}

func NewXoloClientHA(urls []string) *XoloClientHA {
	available := []time.Time{}
	XoloClients := []*XoloClient{}

	for _, url := range urls {
		XoloClients = append(XoloClients, NewXoloClient(url))
		available = append(available, time.Unix(0, 0))
	}
	return &XoloClientHA{
		Rand:        rand.Intn,
		XoloClients: XoloClients,
		available:   available,
	}
}

func (s *XoloClientHA) SendTransaction(ctx context.Context, tx *types.Transaction) (*common.Hash, error) {
	n := len(s.XoloClients)
	i := s.Rand(n)

	for n > 0 {
		if s.available[i].Before(time.Now()) {
			hash, err := s.XoloClients[i].SendTransaction(ctx, tx)
			if err == nil {
				return hash, nil
			}
			log.Errorf("Unable to send transaction to %s : %s", s.XoloClients[i].URL, err)
			s.available[i] = s.available[i].Add(time.Minute)
		}
		i = (i + 1) % len(s.XoloClients)
		n--
	}
	return nil, errors.New("unable to process the transactions in any XoloClient")
}
