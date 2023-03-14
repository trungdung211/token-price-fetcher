package model

import (
	"errors"
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

func (tp *TokenPriceModel) Get(metric string, resolution timeseries.Resolution) (float32, error) {
	for _, ema := range tp.EMA {
		if ema.Metric == metric && ema.Resolution == resolution {
			return ema.Value, nil
		}
	}
	return 0, errors.New("not found in state")
}
