package server

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	echoSwagger "github.com/worldline-go/echo-swagger"
	"github.com/worldline-go/logz/logecho"
	"github.com/worldline-go/tell/metric/metricecho"
	"github.com/worldline-go/wkafka"
	"github.com/ziflex/lecho/v3"

	"github.com/worldline-go/molen/docs"
	"github.com/worldline-go/molen/internal/config"
	"github.com/worldline-go/molen/internal/kafka"
)

type SetConfig struct {
	Client *wkafka.Client
}

// @description github.com/worldline-go/molen
// @BasePath /v1
func Set(ctx context.Context, cfg SetConfig) (*echo.Echo, error) {
	if config.Application.BasePath != "" {
		log.Ctx(ctx).Info().Msgf("base path is set to %s", config.Application.BasePath)
	}

	if err := docs.Info(config.Application.BasePath, config.AppVersion); err != nil {
		return nil, fmt.Errorf("failed to set info; %w", err)
	}

	e := echo.New()

	e.HideBanner = true

	e.Logger = lecho.From(log.With().Str("component", "server").Logger())

	// middlewares
	e.Use(metricecho.HTTPMetrics(nil))
	e.Use(
		middleware.Recover(),
		middleware.CORS(),
	)

	e.Use(
		middleware.RequestID(),
		middleware.RequestLoggerWithConfig(logecho.RequestLoggerConfig()),
		logecho.ZerologLogger(),
	)

	e.Use(
		middleware.Gzip(),
	)

	base := e.Group(config.Application.BasePath)
	base.GET("/swagger/*", echoSwagger.WrapHandler)

	groupV1 := base.Group("/v1")

	// routes"
	kafkaAdmin := cfg.Client.Admin()
	handler := Handler{
		Ctx:            ctx,
		Client:         cfg.Client,
		ProduceMessage: cfg.Client.ProduceRaw,
		ClientAdmin:    kafkaAdmin,
		Group: kafka.Group{
			ClientAdmin: kafkaAdmin,
		},
	}

	groupV1.POST("/publish", handler.Publish)
	groupV1.POST("/group", handler.CreateGroup)

	return e, nil
}
