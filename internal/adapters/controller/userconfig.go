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
	l                 *zap.Logger
}

func NewUserConfigController(us usecase.UserConfig, l *zap.Logger) *UserConfigController {
	return &UserConfigController{us, l}
}

type updateUserConfigReq struct {
	Conditions []string `json:"condition"`
	Tokens     []string `json:"tokens"`
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
		Conditions: conditions,
		Tokens:     req.Tokens,
	}
	u, err := uc.userConfigUsecase.UpdateConfig(
		c.Request.Context(),
		u,
	)
	if err != nil {
		uc.l.Error("userConfigUsecase.UpdateConfig error", zap.Any("err", err))
		c.Error(err)
		return
	}

	res.Data = u
	c.JSON(http.StatusOK, res)
}