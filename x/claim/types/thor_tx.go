package types

type Coin struct {
	Asset  string `json:"asset"`
	Amount string `json:"amount"`
}

type Tx struct {
	ID          string      `json:"id"`
	Chain       string      `json:"chain"`
	FromAddress string      `json:"from_address"`
	ToAddress   string      `json:"to_address"`
	Coins       []Coin      `json:"coins"`
	Gas         interface{} `json:"gas"`
	Memo        string      `json:"memo"`
}

type ObservedTx struct {
	Tx Tx `json:"tx"`
}

type KeysignMetric struct {
	TxID         string      `json:"tx_id"`
	NodeTSSTimes interface{} `json:"node_tss_times"`
}

type ThorChainTxData struct {
	ObservedTx      ObservedTx    `json:"observed_tx"`
	ConsensusHeight int           `json:"consensus_height"`
	FinalisedHeight int           `json:"finalised_height"`
	KeysignMetric   KeysignMetric `json:"keysign_metric"`
}

type ThorTxData struct {
	Hash   string `json:"hash"`
	TxData string `json:"tx_data"`
}
