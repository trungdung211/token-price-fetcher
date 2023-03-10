package usecases

import (
	"context"
	"time"

	"github.com/trungdung211/token-price-fetcher/pkg/timeseries"
)

func (uc *priceUsecase) Load(ctx context.Context) error {
	// load token list
	list, err := uc.userConfigRepo.GetList(ctx)
	if err != nil {
		return err
	}
	tkmap := make(map[string]bool, 0)
	for _, l := range list {
		for _, tk := range l.Tokens {
			tkmap[tk] = true
		}
	}
	tokens := make([]string, 0)
	for tk, _ := range tkmap {
		tokens = append(tokens, tk)
	}

	// init each token
	err = uc.NewToken(ctx, tokens)

	return err
}

func (uc *priceUsecase) loadSeriesFromDB(ctx context.Context, token string) ([]*timeseries.TimeValue, error) {
	out := make([]*timeseries.TimeValue, 0)
	// max resolution is 1 day, so only get price in number of days equal to `capacity`
	limitTime := time.Now().Add(-time.Duration(uc.capacity*24) * time.Hour)
	prices, err := uc.priceRepo.GetLastSeries(ctx, token, uc.supportedResolutions, limitTime, uc.capacity)

	if err != nil {
		return out, err
	}

	for _, p := range prices {
		out = append(out, &timeseries.TimeValue{
			Time:  p.Time,
			Value: p.PriceUSD,
		})
	}

	return out, err
}
