package repo

import (
	"context"
	"errors"
	"time"

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

func (p *priceRepo) GetLastSeries(ctx context.Context, token string, resolutions []model.Resolution, limitTime time.Time, capacity int) ([]*model.Price, error) {
	out := make([]*model.Price, 0)

	err := p.db.NewRaw(
		`select * from (
			SELECT price.*, 
				  rank() OVER (
					  PARTITION BY resolution
					  ORDER BY time DESC
				  )
				FROM price
				WHERE price.token = ?
				  AND time >= ?
				  AND resolution IN (?)
			) price_filter 
		WHERE price_filter.rank <= ?
		ORDER BY price_filter.time ASC;`,
		token, limitTime, resolutions, capacity,
	).Scan(ctx, &out)

	return out, err
}
