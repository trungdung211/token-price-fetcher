package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/trungdung211/token-price-fetcher/internal/entities/model"
	usecase "github.com/trungdung211/token-price-fetcher/internal/usecases/usecases"
	"github.com/trungdung211/token-price-fetcher/pkg/request"
	"go.uber.org/zap"
)

type UserConfigController struct {
	userConfigUsecase usecase.UserConfig
	priceUsecase      usecase.PriceUc
	l                 *zap.Logger
}

func NewUserConfigController(us usecase.UserConfig, pc usecase.PriceUc, l *zap.Logger) *UserConfigController {
	return &UserConfigController{
		userConfigUsecase: us,
		priceUsecase:      pc,
		l:                 l,
	}
}

type updateUserConfigReq struct {
	Conditions     []string `json:"condition"`
	Tokens         []string `json:"tokens"`
	SendNotify     bool     `json:"send_notify"`
	DiscordWebhook *string  `json:"discord_webhook"`
}

// @Summary Update User Config
// @Tags UserConfig
// @Accept json
// @Produce json
// @Param request body updateUserConfigReq true "userconfig create body"
// @Success 200 {object} request.Response{data=model.UserConfig}
// @Router /userconfig/v1/update [post]
// @Security ApiKeyAuth
func (uc *UserConfigController) UpdateConfig(c *gin.Context) {
	res := request.NewResponse()

	// validate
	var req updateUserConfigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		uc.l.Error("ShouldBindJSON error", zap.Any("err", err), zap.Any("req", req))
		c.Error(request.NewError(http.StatusBadRequest, "400", "invalid request body"))
		return
	}

	conditions := make([]model.Condition, 0)
	for _, cond := range req.Conditions {
		if condd, err := model.ParseCondition(cond); err == nil {
			conditions = append(conditions, *condd)
		} else {
			uc.l.Error("Parse condition error", zap.Any("err", err), zap.Any("req", req))
			c.Error(request.NewError(http.StatusBadRequest, "400", "invalid request body"))
			return
		}
	}

	// do update config
	u := &model.UserConfig{
		Conditions:     conditions,
		Tokens:         req.Tokens,
		SendNotify:     req.SendNotify,
		DiscordWebhook: req.DiscordWebhook,
	}
	u, err := uc.userConfigUsecase.UpdateConfig(
		c.Request.Context(),
		u,
	)
	if err != nil {
		uc.l.Error("userConfigUsecase.UpdateConfig error", zap.Any("err", err))
		c.Error(request.NewError(http.StatusBadRequest, "400", err.Error()))
		return
	}

	// update token to price fetcher
	uc.priceUsecase.NewToken(c.Request.Context(), req.Tokens)

	res.Data = u
	c.JSON(http.StatusOK, res)
}
