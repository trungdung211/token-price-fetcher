package repository

import (
	"context"

	"github.com/google/uuid"
	model "github.com/trungdung211/token-price-fetcher/internal/entities/model"
)

type UserConfigRepo interface {
	// FindOne(id uuid.UUID) (*model.UserConfig, error)
	GetByUserId(ctx context.Context, userId uuid.UUID) (*model.UserConfig, error)
	Create(ctx context.Context, u *model.UserConfig) (*model.UserConfig, error)
	Update(ctx context.Context, u *model.UserConfig) (*model.UserConfig, error)
}
