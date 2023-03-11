package external

import (
	"github.com/trungdung211/token-price-fetcher/internal/usecases/external"
	"go.uber.org/zap"
)

type coinGeckoFetcher struct {
	l *zap.Logger
}

func NewCoinGeckoFetcher(l *zap.Logger) external.PriceFetcher {
	return &coinGeckoFetcher{
		l: l,
	}
}

func (c *coinGeckoFetcher) GetPriceUSD(token string) (float32, error) {
	c.l.Debug("Fetch token", zap.Any("token", token))
	return 0.6, nil
}
