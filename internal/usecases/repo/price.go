package repository

import (
	"context"

	"github.com/trungdung211/token-price-fetcher/internal/entities/model"
)

type PriceRepo interface {
	GetPrice(ctx context.Context, token string) (*model.Price, error)
	InsertPrice(ctx context.Context, p *model.Price) (*model.Price, error)
	GetLastSeries(ctx context.Context, token string, durations []model.Resolution, capacity int) ([]*model.Price, error)
}
