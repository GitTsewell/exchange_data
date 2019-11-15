package model

import "time"

type OkexDepth struct {
	Table string `json:"table"`
	Data  []struct {
		Asks         [][]interface{} `json:"asks"`
		Bids         [][]interface{} `json:"bids"`
		InstrumentID string     `json:"instrument_id"`
		Timestamp    time.Time  `json:"timestamp"`
	} `json:"data"`
}
