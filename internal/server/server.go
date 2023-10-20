package server

import (
	"context"
	"fmt"
	"path"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/worldline-go/auth"
	"github.com/worldline-go/auth/jwks"
	"github.com/worldline-go/auth/pkg/authecho"
	echoSwagger "github.com/worldline-go/echo-swagger"
	"github.com/worldline-go/logz/logecho"
	"github.com/worldline-go/tell/metric/metricecho"
	"github.com/worldline-go/wkafka"
	"github.com/ziflex/lecho/v3"

	"github.com/worldline-go/molen/docs"
	"github.com/worldline-go/molen/internal/config"
)

type SetConfig struct {
	Client   *wkafka.Client
	Provider auth.InfProviderExtra
}

// @description github.com/worldline-go/molen
// @description Authorization as "Bearer TOKEN" or use oauth2 login
// @BasePath /v1
// @securityDefinitions.apikey	ApiKeyAuth
// @in header
// @name Authorization
// @securityDefinitions.apikey	ApiKeyAuth
// @in header
// @name Authorization
// @securitydefinitions.oauth2.accessCode	OAuth2AccessCode
// @tokenUrl								[[ .Custom.tokenUrl ]]
// @authorizationUrl						[[ .Custom.authUrl ]]
func Set(ctx context.Context, cfg SetConfig) (*echo.Echo, error) {
	j, err := cfg.Provider.JWTKeyFunc(jwks.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to get key func; %w", err)
	}

	if config.Application.BasePath != "" {
		log.Ctx(ctx).Info().Msgf("base path is set to %s", config.Application.BasePath)
	}

	basePath := path.Join("/", config.Application.BasePath, "/v1")
	if err := docs.Info(basePath, config.AppVersion, cfg.Provider); err != nil {
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

	groupV1 := e.Group(basePath)

	groupV1.GET("/swagger/*", echoSwagger.EchoWrapHandler(func(c *echoSwagger.Config) {
		c.OAuth = &echoSwagger.OAuthConfig{
			ClientId: cfg.Provider.GetClientIDExternal(),
		}
	}))

	producerMessage, err := wkafka.NewProducer(cfg.Client, wkafka.ProducerConfig[Message]{})
	if err != nil {
		return nil, fmt.Errorf("failed to init producer; %w", err)
	}

	// routes"
	handler := Handler{
		Client:         cfg.Client,
		ProduceMessage: producerMessage.Produce,
	}

	mJWT := authecho.MiddlewareJWT(
		authecho.WithKeyFunc(j.Keyfunc),
	)

	groupV1.POST("/publish", handler.Publish, mJWT, authecho.MiddlewareRole(authecho.WithRoles(config.Application.Roles.Write...)))

	return e, nil
}
