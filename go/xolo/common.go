package xolo

type xoloTx struct {
	To    string
	Data  string
	Value string
	Nonce int64
}

type XoloSendTxResult struct {
	Success bool    `json:"success"`
	Error   string  `json:"error,omitempty"`
	TxHash  *string `json:"txhash"`
}

type XoloInfoResult struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
	RpcUrl  string `json:"rpcurl"`
	Address string `json:"address"`
}

type RandNFunc func(int) int
