package usecases

import (
	"errors"
	"fmt"
	"time"

	"github.com/trungdung211/token-price-fetcher/internal/entities/model"
	"github.com/trungdung211/token-price-fetcher/pkg/timeseries"
)

type state struct {
	price float32
	data  map[string]*model.TokenPriceEma
	time  time.Time
}

func newState() *state {
	return &state{
		data: make(map[string]*model.TokenPriceEma, 0),
	}
}

func (s *state) Save(metric string, resolution timeseries.Resolution, value float32) {
	key := fmt.Sprintf("%v:%v", metric, resolution)
	if _, found := s.data[key]; found {
		s.data[key].Value = value
	} else {
		s.data[key] = &model.TokenPriceEma{
			Resolution: resolution,
			Metric:     metric,
			Value:      value,
		}
	}

	// always save the last price
	s.price = value
	s.time = time.Now()
}

func (s *state) Get(metric string, resolution timeseries.Resolution) (float32, error) {
	key := fmt.Sprintf("%v:%v", metric, resolution)
	if val, found := s.data[key]; found {
		return val.Value, nil
	} else {
		return 0, errors.New("not found state")
	}
}

func (s *state) GetAsTokenPriceModel() (*model.TokenPriceModel, error) {
	// get slice of token price
	emas := make([]*model.TokenPriceEma, 0)
	for _, val := range s.data {
		emas = append(emas, val)
	}
	return &model.TokenPriceModel{
		PriceUSD: s.price,
		Time:     s.time,
		EMA:      emas,
	}, nil
}

type priceAgg struct {
	tokenSeries          map[string]*timeseries.MultiResolutionTimeSeries
	tokenPriceState      map[string]*state
	capacity             int
	supportedResolutions []timeseries.Resolution
	insertChan           *chan (*timeseries.TimeValueResolution)
}

func NewPriceAgg(supportedResolutions []timeseries.Resolution, insertChan *chan (*timeseries.TimeValueResolution), capacity int) *priceAgg {
	return &priceAgg{
		tokenSeries:          make(map[string]*timeseries.MultiResolutionTimeSeries, 0),
		tokenPriceState:      make(map[string]*state, 0),
		capacity:             capacity,
		supportedResolutions: supportedResolutions,
		insertChan:           insertChan,
	}
}

func (p *priceAgg) HasToken(token string) bool {
	_, found := p.tokenSeries[token]
	return found
}

func (p *priceAgg) NewToken(token string, series []*timeseries.TimeValue) (err error) {
	// init price series
	ts := timeseries.NewMultiResolutionTimeSeries(
		token,
		p.supportedResolutions,
		p.capacity,
		p.insertChan,
	)
	ts.ReplaceSeries(series)
	p.tokenSeries[token] = ts

	// init state
	p.tokenPriceState[token] = newState()

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
	return m.GetAsTokenPriceModel()
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
	state.Save("ema-7", resolution, ema7)

	// ema20
	ema20, err := timeseries.CalcEMAFromTimeSeries(priceSeries, 20, emaSmooth)
	state.Save("ema-20", resolution, ema20)

	return nil
}
