package config

import (
	"context"
	"fmt"

	"github.com/worldline-go/auth"
	"github.com/worldline-go/igconfig"
	"github.com/worldline-go/tell"
	"github.com/worldline-go/wkafka"
)

var (
	AppName     = "molen"
	AppVersion  = "v0.0.0"
	BuildCommit = "-"
	BuildDate   = "-"
)

var Application = struct {
	Host     string `cfg:"host"`
	Port     string `cfg:"port"`
	LogLevel string `cfg:"log_level"`

	BasePath string `cfg:"base_path"`

	KafkaConfig wkafka.Config `cfg:"kafka_config"`
	// Telemetry configurations
	Telemetry tell.Config

	// AuthService for authentication to API
	AuthService auth.Provider `cfg:"auth_service"`
	Roles       Roles         `cfg:"roles"`
}{
	Host:     "0.0.0.0",
	Port:     "8080",
	LogLevel: "info",
	AuthService: auth.Provider{
		Active: "noop",
	},
}

type Roles struct {
	Admin []string `cfg:"admin"`
	Write []string `cfg:"write"`
}

func Load(ctx context.Context) error {
	if err := igconfig.LoadConfigWithContext(ctx, AppName, &Application); err != nil {
		return fmt.Errorf("unable to load config err: %w", err)
	}

	return nil
}
