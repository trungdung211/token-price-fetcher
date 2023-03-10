package repo

import (
	"context"

	"github.com/google/uuid"
	model "github.com/trungdung211/token-price-fetcher/internal/entities/model"
	repo "github.com/trungdung211/token-price-fetcher/internal/usecases/repo"
	"github.com/uptrace/bun"
)

type userConfigRepo struct {
	db *bun.DB
}

func NewUserConfigRepo(db *bun.DB) repo.UserConfigRepo {
	return &userConfigRepo{db}
}

func (ur *userConfigRepo) GetByUserId(ctx context.Context, userId uuid.UUID) (*model.UserConfig, error) {
	out := &model.UserConfig{}
	err := ur.db.NewSelect().Model(out).Where("user_id = ?", userId).Scan(ctx)
	return out, err
}

func (ur *userConfigRepo) Create(ctx context.Context, u *model.UserConfig) (*model.UserConfig, error) {
	_, err := ur.db.NewInsert().Model(u).Exec(ctx)
	return u, err
}

func (ur *userConfigRepo) Update(ctx context.Context, u *model.UserConfig) (*model.UserConfig, error) {
	_, err := ur.db.NewUpdate().Model(u).WherePK().Exec(ctx)
	return u, err
}

func (ur *userConfigRepo) GetList(ctx context.Context) ([]*model.UserConfig, error) {
	out := make([]*model.UserConfig, 0)
	err := ur.db.NewSelect().Model(&out).Scan(ctx)
	return out, err
}
