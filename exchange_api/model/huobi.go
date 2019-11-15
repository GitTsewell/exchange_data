package model

type HuobiDepth struct {
	Ch   string `json:"ch"`
	Ts   int64  `json:"ts"`
	Tick struct {
		Bids    [][]float64 `json:"bids"`
		Asks    [][]float64 `json:"asks"`
		Version int64       `json:"version"`
		Ts      int64       `json:"ts"`
	} `json:"tick"`
}
