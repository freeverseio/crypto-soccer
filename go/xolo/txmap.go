package xolo

import (
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type txentry struct {
	txhash   common.Hash
	caducity time.Time
}

type TxMap struct {
	sync.Mutex
	cacheTime        time.Duration
	nextCacheCleanup time.Time
	txmap            map[common.Hash]txentry
}

func NewTxMap() *TxMap {
	return &TxMap{
		cacheTime:        time.Second * 60,
		nextCacheCleanup: time.Now(),
		txmap:            map[common.Hash]txentry{},
	}
}

func (m *TxMap) cleanCache() {
	now := time.Now()
	if now.After(m.nextCacheCleanup) {
		for k, v := range m.txmap {
			if v.caducity.After(now) {
				delete(m.txmap, k)
			}
		}
		m.nextCacheCleanup = time.Now().Add(m.cacheTime)
	}
}

func (m *TxMap) Lookup(txhash common.Hash) *common.Hash {
	m.Lock()
	defer m.Unlock()

	if tx, ok := m.txmap[txhash]; ok {
		return &tx.txhash
	} else {
		return nil
	}
}

func (m *TxMap) Add(txhashint common.Hash, txhashreal common.Hash) {
	m.Lock()
	defer m.Unlock()

	m.cleanCache()

	m.txmap[txhashint] = txentry{
		txhash:   txhashreal,
		caducity: time.Now().Add(m.cacheTime),
	}
}
