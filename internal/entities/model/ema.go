package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Resolution int

const (
	EMA_RESOLUT_1_MIN Resolution = iota
	EMA_RESOLUT_1_HOUR
	EMA_RESOLUT_4_HOUR
	EMA_RESOLUT_1_DAY
)

func (r *Resolution) ToString() string {
	switch *r {
	case EMA_RESOLUT_1_MIN:
		return "1min"
	case EMA_RESOLUT_1_HOUR:
		return "1hour"
	case EMA_RESOLUT_4_HOUR:
		return "4hours"
	case EMA_RESOLUT_1_DAY:
		return "1day"
	}
	return ""
}

type Ema struct {
	bun.BaseModel `bun:"table:ema" swaggerignore:"true"`

	Id         uuid.UUID  `bun:",pk,type:uuid"`
	Time       time.Time  `bun:"time,type:timestamptz"`
	Token      string     `bun:"token,type:varchar(24)"`
	Resolution Resolution `bun:"resolution,type:int"`
	Value      float32    `bun:"value,type:real"`
}

var _ bun.BeforeAppendModelHook = (*Ema)(nil)

func (e *Ema) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		if e.Id == uuid.Nil {
			e.Id = uuid.New()
		}
	}
	return nil
}
