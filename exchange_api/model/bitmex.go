package model

import "time"

type BitmexDepth struct {
	Table  string `json:"table"`
	Action string `json:"action"`
	Data   []struct {
		Symbol    string      `json:"symbol"`
		Asks      [][]float64     `json:"asks"`
		Timestamp time.Time   `json:"timestamp"`
		Bids      [][]float64 `json:"bids"`
	} `json:"data"`
}
