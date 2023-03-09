package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	dbpkg "github.com/trungdung211/token-price-fetcher/pkg/postgres"

	// Swagger docs.
	_ "github.com/trungdung211/token-price-fetcher/gen/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Run
// Swagger spec:
// @title       Go API
// @description Project swagger
// @version     1.0
// @host        http://localhost:8080
// @BasePath    /api/v1
// @schemes http
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func Run() {
	// init logger
	l := newLogger()
	// // init db
	_, err := dbpkg.NewPostgresDb(viper.GetString("postgres.uri"), false)
	if err != nil {
		panic(err)
	}

	// HTTP Server
	handler := gin.New()

	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())
	handler.Use(errorHandler())

	// Swagger
	swaggerHandler := ginSwagger.WrapHandler(swaggerFiles.Handler)
	handler.GET("/swagger/*any", swaggerHandler)

	httpServer := &http.Server{
		Handler:      handler,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		Addr:         viper.GetString("address"),
	}

	// Run server
	errChan := make(chan error, 1)
	go func() {
		errChan <- httpServer.ListenAndServe()
		close(errChan)
	}()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err := <-errChan:
		l.Error("app - Run - httpServer.Notify error", zap.Any("err", err))
	}

	// gracefully shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	httpServer.Shutdown(ctx)
}

func newLogger() *zap.Logger {
	env := viper.GetString("env")
	var logger *zap.Logger
	switch env {
	case "development":
		logger, _ = zap.NewDevelopment()
	default:
		logger, _ = zap.NewProduction(zap.IncreaseLevel(zap.WarnLevel))
		gin.SetMode(gin.ReleaseMode)
	}
	return logger
}