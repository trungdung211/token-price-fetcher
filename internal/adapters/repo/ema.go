package repo

import (
	"context"

	model "github.com/trungdung211/token-price-fetcher/internal/entities/model"
	repo "github.com/trungdung211/token-price-fetcher/internal/usecases/repo"
	"github.com/uptrace/bun"
)

type emaRepo struct {
	db *bun.DB
}

func NewEmaRepo(db *bun.DB) repo.EmaRepo {
	return &emaRepo{db}
}

func (r *emaRepo) InsertMany(ctx context.Context, ema []*model.Ema) error {
	_, err := r.db.NewInsert().Model(&ema).Exec(ctx)
	return err
}
