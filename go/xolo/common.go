package xolo

type xoloTx struct {
	TxID    string
	ChainID int64
	Pool    string
	To      string
	Data    string
	Value   string
	Nonce   int64
}

type XoloApiBaseResult struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

type XoloApiTranslateTxResult struct {
	Success bool    `json:"success"`
	Error   string  `json:"error,omitempty"`
	Tx      *string `json:"tx"`
}
