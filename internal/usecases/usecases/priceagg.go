package usercases

import (
	"errors"

	"github.com/spf13/viper"
	"github.com/trungdung211/token-price-fetcher/pkg/timeseries"
)

type priceAgg struct {
	tokenSeries map[string]*timeseries.MultiResolutionTimeSeries
	capacity    int
	insertChan  *chan (*timeseries.TimeValueResolution)
}

func NewPriceAgg(insertChan *chan (*timeseries.TimeValueResolution)) *priceAgg {
	return &priceAgg{
		tokenSeries: make(map[string]*timeseries.MultiResolutionTimeSeries, 0),
		capacity:    viper.GetInt("price.capacity"),
		insertChan:  insertChan,
	}
}

func (p *priceAgg) HasToken(token string) bool {
	_, found := p.tokenSeries[token]
	return found
}

func (p *priceAgg) NewToken(token string, series []*timeseries.TimeValue) (err error) {
	RESOLUTIONS := []timeseries.Resolution{
		timeseries.TIME_RESOLUTION_1_MIN,
		timeseries.TIME_RESOLUTION_1_HOUR,
		timeseries.TIME_RESOLUTION_4_HOURS,
		timeseries.TIME_RESOLUTION_1_DAY,
	}

	ts := timeseries.NewMultiResolutionTimeSeries(
		token,
		RESOLUTIONS,
		p.capacity,
		p.insertChan,
	)
	ts.ReplaceSeries(series)
	p.tokenSeries[token] = ts

	return
}

func (p *priceAgg) GetToken(token string) (*timeseries.MultiResolutionTimeSeries, error) {
	ts, found := p.tokenSeries[token]
	if !found {
		return nil, errors.New("not found token")
	}
	return ts, nil
}
