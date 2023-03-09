package usercases

import (
	"context"

	model "github.com/trungdung211/token-price-fetcher/internal/entities/model"
	repo "github.com/trungdung211/token-price-fetcher/internal/usecases/repo"
)

type priceUsecase struct {
	priceRepo repo.PriceRepo
	emaRepo   repo.EmaRepo
}

type PriceUc interface {
	NewToken(ctx context.Context, tokens []string) error
	GetTokenPrice(ctx context.Context, token string) (*model.TokenPriceModel, error)
	FetchForever() error
}

func NewPriceUsecase(priceRepo repo.PriceRepo, emaRepo repo.EmaRepo) PriceUc {
	return &priceUsecase{
		priceRepo: priceRepo,
		emaRepo:   emaRepo,
	}
}

func (uc *priceUsecase) NewToken(ctx context.Context, tokens []string) error {
	return nil
}

func (uc *priceUsecase) GetTokenPrice(ctx context.Context, token string) (*model.TokenPriceModel, error) {
	return nil, nil
}

func (uc *priceUsecase) FetchForever() error {
	return nil
}
