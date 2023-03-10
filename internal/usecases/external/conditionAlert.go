package external

import "github.com/trungdung211/token-price-fetcher/internal/entities/model"

type ConditionAlert interface {
	Alert(config *model.UserConfig, message string) error
}
