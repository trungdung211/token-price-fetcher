package repo

import (
	"context"
	"errors"

	model "github.com/trungdung211/token-price-fetcher/internal/entities/model"
	repo "github.com/trungdung211/token-price-fetcher/internal/usecases/repo"
	"github.com/uptrace/bun"
)

type priceRepo struct {
	db *bun.DB
}

func NewPriceRepo(db *bun.DB) repo.PriceRepo {
	return &priceRepo{db}
}

func (p *priceRepo) GetPrice(ctx context.Context, token string) (*model.Price, error) {
	out := []model.Price{}
	err := p.db.NewSelect().
		Model(out).
		Where("token = ?", token).
		Order("time DESC").
		Limit(1).Scan(ctx)
	if err != nil {
		return nil, err
	}
	if len(out) == 0 {
		return nil, errors.New("no fetched price")
	}

	m := out[0]
	return &m, nil
}

func (p *priceRepo) InsertPrice(ctx context.Context, m *model.Price) (*model.Price, error) {
	_, err := p.db.NewInsert().Model(m).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return m, err
}
