package model

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Condition string

const (
	CONDITION_DIPS_1H_EMA_20 = "DIPS_1H_EMA_20"
)

func ParseCondition(s string) (*Condition, error) {
	m := map[string]Condition{
		"DIPS_1H_EMA_20": CONDITION_DIPS_1H_EMA_20,
	}
	if c, found := m[s]; found {
		return &c, nil
	} else {
		return nil, errors.New("not found condition")
	}
}

type UserConfig struct {
	bun.BaseModel `bun:"table:user_config" swaggerignore:"true"`

	Id     uuid.UUID `bun:",pk,type:uuid" json:"id"`
	UserId uuid.UUID `bun:"user_id,type:uuid" json:"user_id"`
	// discord info

	// condition
	Conditions []Condition `bun:"conditions,type:jsonb" json:"conditions"`

	// token list
	Tokens []string `bun:"tokens,type:jsonb" json:"token"`

	// audit
	CreatedAt time.Time `bun:"created_at,type:timestamptz" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,type:timestamptz" json:"updated_at"`
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
