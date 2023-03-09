package model

import "time"

type TokenPriceModel struct {
	PriceUSD float32   `json:"price_usd"`
	EMA7_1M  float32   `json:"ema7_1m"`
	EMA7_1H  float32   `json:"ema7_1h"`
	EMA7_4H  float32   `json:"ema7_4h"`
	EMA7_1D  float32   `json:"ema7_1d"`
	EMA20_1M float32   `json:"ema20_1m"`
	EMA20_1H float32   `json:"ema20_1h"`
	EMA20_4H float32   `json:"ema20_4h"`
	EMA20_1D float32   `json:"ema20_1d"`
	Time     time.Time `json:"time"`
}
