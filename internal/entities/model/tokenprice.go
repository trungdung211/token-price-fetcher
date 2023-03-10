package model

import (
	"time"

	"github.com/trungdung211/token-price-fetcher/pkg/timeseries"
)

type TokenPriceEma struct {
	Metric     string
	Resolution timeseries.Resolution
	Value      float32
}

type TokenPriceModel struct {
	PriceUSD float32          `json:"price_usd"`
	EMA      []*TokenPriceEma `json:"ema"`
	Time     time.Time        `json:"time"`
}
