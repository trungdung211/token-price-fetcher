package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Condition int

const (
	CONDITION_DIPS_1H_EMA_20 = iota
)

type UserConfig struct {
	bun.BaseModel `bun:"table:user_config" swaggerignore:"true"`

	Id     uuid.UUID `bun:",pk,type:uuid"`
	UserId uuid.UUID `bun:"user_id,type:uuid"`
	// discord info

	// condition
	Conditions []Condition `bun:"conditions,type:jsonb"`

	// token list
	Tokens []string `bun:"tokens,type:jsonb"`

	// audit
	CreatedAt time.Time `bun:"created_at,type:timestamptz"`
	UpdatedAt time.Time `bun:"updated_at,type:timestamptz"`
}

var _ bun.BeforeAppendModelHook = (*UserConfig)(nil)

func (e *UserConfig) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		if e.Id == uuid.Nil {
			e.Id = uuid.New()
		}

		e.CreatedAt = time.Now()
		e.UpdatedAt = time.Now()
	case *bun.UpdateQuery:
		e.UpdatedAt = time.Now()
	}
	return nil
}
