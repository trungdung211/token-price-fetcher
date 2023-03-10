package usecases

import (
	"fmt"

	"github.com/trungdung211/token-price-fetcher/internal/entities/model"
	"github.com/trungdung211/token-price-fetcher/pkg/timeseries"
)

type Condition interface {
	Condition() model.Condition
	IsMatched(state *model.TokenPriceModel) (string, bool)
}

type ConditionDipsEMA20R1H struct {
}

func (c *ConditionDipsEMA20R1H) Condition() model.Condition {
	return model.CONDITION_DIPS_1H_EMA_20
}

func (c *ConditionDipsEMA20R1H) IsMatched(state *model.TokenPriceModel) (out string, matched bool) {
	if ema, err := state.Get("ema-20", timeseries.TIME_RESOLUTION_1_HOUR); err == nil {
		matched = state.PriceUSD < ema && ema > 0
		if matched {
			out = fmt.Sprintf(`The price dips under its one-hour EMA-20 (price_usd = %v, ema = %v)`, state.PriceUSD, ema)
		}
	}
	return
}

type ConditionDipsEMA7R4H struct {
}

func (c *ConditionDipsEMA7R4H) Condition() model.Condition {
	return model.CONDITION_DIPS_4H_EMA_7
}

func (c *ConditionDipsEMA7R4H) IsMatched(state *model.TokenPriceModel) (out string, matched bool) {
	if ema, err := state.Get("ema-7", timeseries.TIME_RESOLUTION_4_HOURS); err == nil {
		matched = state.PriceUSD < ema && ema > 0
		if matched {
			out = fmt.Sprintf(`The price dips under its four-hours EMA-7 (price_usd = %v, ema = %v)`, state.PriceUSD, ema)
		}
	}
	return
}

type ConditionDipsEMA7R1M struct {
}

func (c *ConditionDipsEMA7R1M) Condition() model.Condition {
	return model.CONDITION_DIPS_1M_EMA_7
}

func (c *ConditionDipsEMA7R1M) IsMatched(state *model.TokenPriceModel) (out string, matched bool) {
	if ema, err := state.Get("ema-7", timeseries.TIME_RESOLUTION_1_MIN); err == nil {
		matched = state.PriceUSD < ema && ema > 0
		if matched {
			out = fmt.Sprintf(`The price dips under its one-minute EMA-7 (price_usd = %v, ema = %v)`, state.PriceUSD, ema)
		}
	}
	return
}

func NewCondition(condition model.Condition) Condition {
	switch condition {
	case model.CONDITION_DIPS_1H_EMA_20:
		return &ConditionDipsEMA20R1H{}
	case model.CONDITION_DIPS_4H_EMA_7:
		return &ConditionDipsEMA7R4H{}
	case model.CONDITION_DIPS_1M_EMA_7:
		return &ConditionDipsEMA7R1M{}
	}
	return &ConditionDipsEMA7R1M{}
}
