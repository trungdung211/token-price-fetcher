package usercases

import (
	"context"

	"github.com/spf13/viper"
	model "github.com/trungdung211/token-price-fetcher/internal/entities/model"
	"github.com/trungdung211/token-price-fetcher/internal/usecases/external"
	repo "github.com/trungdung211/token-price-fetcher/internal/usecases/repo"
	"github.com/trungdung211/token-price-fetcher/pkg/timeseries"
)

type priceUsecase struct {
	priceRepo    repo.PriceRepo
	emaRepo      repo.EmaRepo
	priceFetcher external.PriceFetcher
	priceAgg     *priceAgg
	insertChan   chan (*timeseries.TimeValueResolution)
}

type PriceUc interface {
	NewToken(ctx context.Context, tokens []string) error
	GetTokenPrice(ctx context.Context, token string) (*model.TokenPriceModel, error)
	FetchForever() error
}

func NewPriceUsecase(priceRepo repo.PriceRepo, emaRepo repo.EmaRepo, priceFetcher external.PriceFetcher) PriceUc {
	insertChan := make(chan (*timeseries.TimeValueResolution), viper.GetInt("price.emachan_capacity"))
	return &priceUsecase{
		priceRepo:    priceRepo,
		emaRepo:      emaRepo,
		priceFetcher: priceFetcher,
		priceAgg:     NewPriceAgg(&insertChan),
		insertChan:   insertChan,
	}
}

func (uc *priceUsecase) NewToken(ctx context.Context, tokens []string) error {
	for _, token := range tokens {
		// check exists token
		if uc.priceAgg.HasToken(token) {
			continue
		}
		// TODO : get old price data from db
		// do new token
		uc.priceAgg.NewToken(token, []*timeseries.TimeValue{})
	}

	return nil
}

func (uc *priceUsecase) GetTokenPrice(ctx context.Context, token string) (*model.TokenPriceModel, error) {
	return nil, nil
}

func (uc *priceUsecase) FetchForever() error {
	return nil
}
