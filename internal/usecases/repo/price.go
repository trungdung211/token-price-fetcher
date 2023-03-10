package repository

import (
	"context"
	"time"

	"github.com/trungdung211/token-price-fetcher/internal/entities/model"
	"github.com/trungdung211/token-price-fetcher/pkg/timeseries"
)

type PriceRepo interface {
	GetPrice(ctx context.Context, token string) (*model.Price, error)
	InsertPrice(ctx context.Context, p *model.Price) (*model.Price, error)
	GetLastSeries(ctx context.Context, token string, durations []timeseries.Resolution, limitTime time.Time, capacity int) ([]*model.Price, error)
}
