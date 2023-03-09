package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Price struct {
	bun.BaseModel `bun:"table:price" swaggerignore:"true"`

	Id       uuid.UUID `bun:",pk,type:uuid"`
	Time     time.Time `bun:"time,type:timestamptz"`
	Token    string    `bun:"token,type:varchar(24)"`
	PriceUSD float32   `bun:"price_usd,type:real"`
}

var _ bun.BeforeAppendModelHook = (*Price)(nil)

func (e *Price) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		if e.Id == uuid.Nil {
			e.Id = uuid.New()
		}
	}
	return nil
}
