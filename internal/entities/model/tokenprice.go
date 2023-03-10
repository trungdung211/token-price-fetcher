package model

import "time"

type TokenPriceEma struct {
	Metric     string
	Resolution Resolution
	Value      float32
}

type TokenPriceModel struct {
	PriceUSD float32          `json:"price_usd"`
	EMA      []*TokenPriceEma `json:"ema"`
	Time     time.Time        `json:"time"`
}
