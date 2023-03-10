package usecases

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
	l                    *zap.Logger
	priceRepo            repo.PriceRepo
	userConfigRepo       repo.UserConfigRepo
	priceFetcher         external.PriceFetcher
	priceAgg             *priceAgg
	insertChan           chan (*timeseries.TimeValueResolution)
	emaSmooth            float32
	capacity             int
	supportedResolutions []timeseries.Resolution
	trigger              TriggerCondition
}

type PriceUc interface {
	Load(ctx context.Context) error
	NewToken(ctx context.Context, tokens []string) error
	GetTokenPrice(ctx context.Context, token string) (*model.TokenPriceModel, error)
	FetchForever() error
}

func NewPriceUsecase(l *zap.Logger, priceRepo repo.PriceRepo, userConfigRepo repo.UserConfigRepo, priceFetcher external.PriceFetcher, trigger TriggerCondition) PriceUc {
	insertChan := make(chan (*timeseries.TimeValueResolution), viper.GetInt("price.emachan_capacity"))
	emaSmooth := (float32)(viper.GetFloat64("price.ema_smooth"))
	capacity := viper.GetInt("price.capacity")
	supportedResolutions := []timeseries.Resolution{
		timeseries.TIME_RESOLUTION_1_MIN,
		timeseries.TIME_RESOLUTION_1_HOUR,
		timeseries.TIME_RESOLUTION_4_HOURS,
		timeseries.TIME_RESOLUTION_1_DAY,
	}

	return &priceUsecase{
		l:                    l,
		priceRepo:            priceRepo,
		userConfigRepo:       userConfigRepo,
		priceFetcher:         priceFetcher,
		priceAgg:             NewPriceAgg(supportedResolutions, &insertChan, capacity),
		trigger:              trigger,
		insertChan:           insertChan,
		emaSmooth:            emaSmooth,
		capacity:             capacity,
		supportedResolutions: supportedResolutions,
	}
}

func (uc *priceUsecase) NewToken(ctx context.Context, tokens []string) error {
	for _, token := range tokens {
		// check exists token
		if uc.priceAgg.HasToken(token) {
			continue
		}
		// TODO : get old price data from db
		series, _ := uc.loadSeriesFromDB(ctx, token)
		uc.l.Info("Load last series for token", zap.Any("count", len(series)), zap.Any("token", token))

		// do new token
		uc.priceAgg.NewToken(token, series)
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
			token := inserted.Name
			// persist to db
			uc.priceRepo.InsertPrice(ctx, &model.Price{
				Time:       inserted.TV.Time,
				PriceUSD:   inserted.TV.Value,
				Token:      token,
				Resolution: inserted.Resolution,
			})
			// calc ema
			uc.priceAgg.CalcEMA(token, inserted.Resolution, inserted.TV.Value, inserted.TV.Time, uc.emaSmooth)

			// trigger condition
			updatedState, _ := uc.priceAgg.GetTokenPriceState(inserted.Name)
			go uc.trigger.Trigger(ctx, token, updatedState)
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
				series.Add(price, time.Now(), true)
			}
		}()
	}

	return nil
}
