package usecases

import (
	"context"
	"fmt"
	"strings"

	"github.com/trungdung211/token-price-fetcher/internal/entities/model"
	"github.com/trungdung211/token-price-fetcher/internal/usecases/external"
	repo "github.com/trungdung211/token-price-fetcher/internal/usecases/repo"
	"go.uber.org/zap"
)

type triggerCondition struct {
	l              *zap.Logger
	userConfigRepo repo.UserConfigRepo
	conditionAlert external.ConditionAlert
}

type TriggerCondition interface {
	Trigger(ctx context.Context, token string, state *model.TokenPriceModel) error
}

func NewTriggerCondition(
	l *zap.Logger,
	userConfigRepo repo.UserConfigRepo,
	conditionAlert external.ConditionAlert) TriggerCondition {
	return &triggerCondition{
		l:              l,
		userConfigRepo: userConfigRepo,
		conditionAlert: conditionAlert,
	}
}

func (tc *triggerCondition) Trigger(ctx context.Context, token string, state *model.TokenPriceModel) error {
	// check state time before send notify

	// get all user
	configs, err := tc.userConfigRepo.GetByToken(ctx, token)
	if err != nil {
		tc.l.Error("userConfigRepo.GetByToken err", zap.Any("err", err))
		return err
	}

	// check all conditions
	for _, c := range configs {
		if !c.SendNotify {
			continue
		}
		for _, cond := range c.Conditions {
			condObj := NewCondition(cond)
			if message, matched := condObj.IsMatched(state); matched {
				tc.l.Debug("Trigger condition", zap.Any("user", c.UserId), zap.Any("message", message))
				// trigger message to user discord
				message = fmt.Sprintf("Token (%v): %s", strings.ToUpper(token), message)
				tc.conditionAlert.Alert(c, message)
			}
		}
	}
	return nil
}
