package model

type BinanceDepth struct {
	Stream string `json:"stream"`
	Data   struct {
		LastUpdateID int        `json:"lastUpdateId"`
		Bids         [][]interface{} `json:"bids"`
		Asks         [][]interface{} `json:"asks"`
	} `json:"data"`
}
