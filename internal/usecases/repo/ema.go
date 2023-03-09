package repository

import (
	"context"

	"github.com/trungdung211/token-price-fetcher/internal/entities/model"
)

type EmaRepo interface {
	InsertMany(ctx context.Context, ema []*model.Ema) error
}
