package external

import "github.com/trungdung211/token-price-fetcher/internal/usecases/external"

type coinGeckoFetcher struct {
}

func NewCoinGeckoFetcher() external.PriceFetcher {
	return &coinGeckoFetcher{}
}

func (c *coinGeckoFetcher) GetPriceUSD(token string) (float32, error) {
	return 1, nil
}
