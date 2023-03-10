package usercases

import (
	"context"
	"time"

	"github.com/spf13/viper"
	model "github.com/trungdung211/token-price-fetcher/internal/entities/model"
	"github.com/trungdung211/token-price-fetcher/internal/usecases/external"
	repo "github.com/trungdung211/token-price-fetcher/internal/usecases/repo"
	"github.com/trungdung211/token-price-fetcher/pkg/timeseries"
	"go.uber.org/zap"
)

type priceUsecase struct {
	l            *zap.Logger
	priceRepo    repo.PriceRepo
	emaRepo      repo.EmaRepo
	priceFetcher external.PriceFetcher
	priceAgg     *priceAgg
	insertChan   chan (*timeseries.TimeValueResolution)
	emaSmooth    float32
}

type PriceUc interface {
	NewToken(ctx context.Context, tokens []string) error
	GetTokenPrice(ctx context.Context, token string) (*model.TokenPriceModel, error)
	FetchForever() error
}

func NewPriceUsecase(l *zap.Logger, priceRepo repo.PriceRepo, emaRepo repo.EmaRepo, priceFetcher external.PriceFetcher) PriceUc {
	insertChan := make(chan (*timeseries.TimeValueResolution), viper.GetInt("price.emachan_capacity"))
	emaSmooth := (float32)(viper.GetFloat64("price.ema_smooth"))

	return &priceUsecase{
		l:            l,
		priceRepo:    priceRepo,
		emaRepo:      emaRepo,
		priceFetcher: priceFetcher,
		priceAgg:     NewPriceAgg(&insertChan),
		insertChan:   insertChan,
		emaSmooth:    emaSmooth,
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
	priceModel, err := uc.priceAgg.GetTokenPriceState(token)
	return priceModel, err
}

func (uc *priceUsecase) FetchForever() error {
	ctx := context.Background()

	// calculate ema and trigger notification if there are conditions is met
	go func() {
		for {
			inserted := <-uc.insertChan
			// persist to db
			uc.priceRepo.InsertPrice(ctx, &model.Price{
				Time:       inserted.TV.Time,
				PriceUSD:   inserted.TV.Value,
				Token:      inserted.Name,
				Resolution: model.Resolution(inserted.Resolution),
			})
			// calc ema
			uc.priceAgg.CalcEMA(inserted.Name, inserted.Resolution, inserted.TV.Value, inserted.TV.Time, uc.emaSmooth)

			// trigger condition
		}
	}()

	tokenPriceFetch := make(chan (string), 100)

	// continuously fetch new price
	go func() {
		interval := viper.GetInt("price.fetch_interval")
		for {
			tokenList, _ := uc.priceAgg.GetTokenList()
			for _, token := range tokenList {
				tokenPriceFetch <- token
			}
			time.Sleep(time.Duration(interval) * time.Second)
		}
	}()

	// fetch price worker
	concurrent := viper.GetInt("price.fetch_concurrent")
	for i := 0; i < concurrent; i++ {
		go func() {
			for {
				token := <-tokenPriceFetch
				price, err := uc.priceFetcher.GetPriceUSD(token)
				if err != nil {
					uc.l.Error("priceFetcher.GetPriceUSD error", zap.Any("err", err))
					continue
				}
				series, err := uc.priceAgg.GetToken(token)
				if err != nil {
					continue
				}
				series.Add(price, time.Now())
			}
		}()
	}

	return nil
}
