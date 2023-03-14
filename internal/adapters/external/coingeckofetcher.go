package external

import (
	"github.com/JulianToledano/goingecko"
	"github.com/trungdung211/token-price-fetcher/internal/usecases/external"
	"go.uber.org/zap"
)

type coinGeckoFetcher struct {
	l        *zap.Logger
	cgClient *goingecko.Client
}

func NewCoinGeckoFetcher(l *zap.Logger) external.PriceFetcher {
	return &coinGeckoFetcher{
		l:        l,
		cgClient: goingecko.NewClient(nil),
	}
}

func (c *coinGeckoFetcher) GetPriceUSD(token string) (float32, error) {
	c.l.Debug("Fetch token", zap.Any("token", token))
	// return 0.6, nil
	data, err := c.cgClient.CoinsId(token, true, true, true, false, false, false)
	if err != nil {
		c.l.Error("cgClient.CoinsId err", zap.Any("token", token), zap.Any("err", err))
		return 0, err
	}
	price := (float32)(data.MarketData.CurrentPrice.Usd)
	return price, nil
}
