package app

import (
	"github.com/gin-gonic/gin"
	controller "github.com/trungdung211/token-price-fetcher/internal/adapters/controller"
	usecases "github.com/trungdung211/token-price-fetcher/internal/usecases/usecases"
	"go.uber.org/zap"
)

func initRouter(handler *gin.Engine, l *zap.Logger, us usecases.UserConfig, pu usecases.PriceUc) {
	// userconfig
	ucc := controller.NewUserConfigController(
		us, pu, l,
	)
	h := handler.Group("/userconfig/v1")
	h.POST("/update", ucc.UpdateConfig)

	// price
	pc := controller.NewPriceController(pu, l)
	h = handler.Group("/prices/v1")
	h.GET("/token/{token}", pc.GetPrice)
}
