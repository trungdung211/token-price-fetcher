package usercases

import (
	"errors"
	"time"

	"github.com/spf13/viper"
	model "github.com/trungdung211/token-price-fetcher/internal/entities/model"
	"github.com/trungdung211/token-price-fetcher/pkg/timeseries"
)

type priceAgg struct {
	tokenSeries     map[string]*timeseries.MultiResolutionTimeSeries
	tokenPriceState map[string]*model.TokenPriceModel
	capacity        int
	insertChan      *chan (*timeseries.TimeValueResolution)
}

func NewPriceAgg(insertChan *chan (*timeseries.TimeValueResolution)) *priceAgg {
	return &priceAgg{
		tokenSeries:     make(map[string]*timeseries.MultiResolutionTimeSeries, 0),
		tokenPriceState: make(map[string]*model.TokenPriceModel, 0),
		capacity:        viper.GetInt("price.capacity"),
		insertChan:      insertChan,
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

	// init price series
	ts := timeseries.NewMultiResolutionTimeSeries(
		token,
		RESOLUTIONS,
		p.capacity,
		p.insertChan,
	)
	ts.ReplaceSeries(series)
	p.tokenSeries[token] = ts

	// init state
	p.tokenPriceState[token] = &model.TokenPriceModel{
		Time: time.Now(),
	}

	return
}

func (p *priceAgg) GetToken(token string) (*timeseries.MultiResolutionTimeSeries, error) {
	ts, found := p.tokenSeries[token]
	if !found {
		return nil, errors.New("not found token")
	}
	return ts, nil
}

func (p *priceAgg) GetTokenList() ([]string, error) {
	keys := make([]string, 0)
	for k := range p.tokenSeries {
		keys = append(keys, k)
	}
	return keys, nil
}

func (p *priceAgg) GetTokenPriceState(token string) (*model.TokenPriceModel, error) {
	m, found := p.tokenPriceState[token]
	if !found {
		return nil, errors.New("token not found")
	}
	return m, nil
}

func (p *priceAgg) CalcEMA(token string, resolution timeseries.Resolution, value float32, ts time.Time, emaSmooth float32) error {
	price, _ := p.GetToken(token)
	priceSeries, err := price.GetSeries(resolution)
	if err != nil {
		return err
	}

	state := p.tokenPriceState[token]
	// ema7
	ema7, err := timeseries.CalcEMAFromTimeSeries(priceSeries, 7, emaSmooth)
	if err != nil {
		switch resolution {
		case model.EMA_RESOLUT_1_MIN:
			state.EMA7_1M = ema7
		case model.EMA_RESOLUT_1_HOUR:
			state.EMA7_1H = ema7
		case model.EMA_RESOLUT_4_HOUR:
			state.EMA7_4H = ema7
		case model.EMA_RESOLUT_1_DAY:
			state.EMA7_1D = ema7
		}
	}

	// ema20
	ema20, err := timeseries.CalcEMAFromTimeSeries(priceSeries, 20, emaSmooth)
	if err != nil {
		switch resolution {
		case model.EMA_RESOLUT_1_MIN:
			state.EMA20_1M = ema20
		case model.EMA_RESOLUT_1_HOUR:
			state.EMA20_1H = ema20
		case model.EMA_RESOLUT_4_HOUR:
			state.EMA20_4H = ema20
		case model.EMA_RESOLUT_1_DAY:
			state.EMA20_1D = ema20
		}
	}

	return nil
}
