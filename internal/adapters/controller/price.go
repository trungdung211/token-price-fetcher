package controller

import (
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	usecases "github.com/trungdung211/token-price-fetcher/internal/usecases/usecases"
	"github.com/trungdung211/token-price-fetcher/pkg/request"
)

type PriceController struct {
	priceUsecase usecases.PriceUc
	l            *zap.Logger
}

func NewPriceController(priceUc usecases.PriceUc, l *zap.Logger) *PriceController {
	return &PriceController{priceUc, l}
}

type priceEma struct {
	Metric     string  `json:"metric"`
	Resolution string  `json:"resolution"`
	Value      float32 `json:"value"`
}

type priceModel struct {
	PriceUSD float32     `json:"price_usd"`
	EMA      []*priceEma `json:"ema"`
	Time     time.Time   `json:"time"`
}

// @Summary Get Token Price
// @Tags Price
// @Accept json
// @Produce json
// @Param token path string true "token"
// @Success 200 {object} request.Response{data=model.TokenPriceModel}
// @Router /prices/v1/token/{token} [get]
// @Security ApiKeyAuth
func (pc *PriceController) GetPrice(c *gin.Context) {
	res := request.NewResponse()

	token := c.Param("token")
	// get info
	data, err := pc.priceUsecase.GetTokenPrice(c.Request.Context(), token)
	if err != nil {
		pc.l.Error("priceUsecase.GetTokenPrice error", zap.Any("err", err), zap.Any("token", token))
		c.Error(request.NewError(http.StatusBadRequest, "404", err.Error()))
		return
	}

	ema := make([]*priceEma, 0)
	for _, m := range data.EMA {
		ema = append(ema, &priceEma{
			Metric:     m.Metric,
			Resolution: m.Resolution.ToString(),
			Value:      m.Value,
		})
	}

	res.Data = &priceModel{
		PriceUSD: data.PriceUSD,
		Time:     data.Time,
		EMA:      ema,
	}
	c.JSON(http.StatusOK, res)
}
