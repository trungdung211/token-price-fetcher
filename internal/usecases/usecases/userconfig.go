package usercases

import (
	"context"

	"github.com/google/uuid"
	"github.com/trungdung211/token-price-fetcher/internal/entities/model"
	repo "github.com/trungdung211/token-price-fetcher/internal/usecases/repo"
)

var (
	DEFAULT_USER_ID uuid.UUID
)

func init() {
	// have not implemented authentication, fake default user_id
	DEFAULT_USER_ID, _ = uuid.Parse("8300d4dc-be93-11ed-960f-de8e0d88801c")
}

type userConfigUsecase struct {
	userRepo repo.UserConfigRepo
}

type UserConfig interface {
	UpdateConfig(ctx context.Context, u *model.UserConfig) (*model.UserConfig, error)
}

func NewUserConfigUsecase(r repo.UserConfigRepo) UserConfig {
	return &userConfigUsecase{r}
}

func (uu *userConfigUsecase) UpdateConfig(ctx context.Context, u *model.UserConfig) (*model.UserConfig, error) {
	// fake userId
	u.UserId = DEFAULT_USER_ID

	data, err := uu.userRepo.GetByUserId(ctx, u.UserId)
	if err != nil {
		data, err = uu.userRepo.Create(ctx, u)
	} else {
		u.Id = data.Id
		data, err = uu.userRepo.Update(ctx, u)
	}
	return data, err
}
