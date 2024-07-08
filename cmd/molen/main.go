package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/worldline-go/initializer"
	"github.com/worldline-go/logz"
	"github.com/worldline-go/tell"
	"github.com/worldline-go/wkafka"

	"github.com/worldline-go/molen/internal/config"
	"github.com/worldline-go/molen/internal/server"
)

var (
	version = "v0.0.0"
	commit  = "-"
	date    = "-"
)

func setBuildInfo() {
	config.AppVersion = version
	config.BuildCommit = commit
	config.BuildDate = date
}

func main() {
	setBuildInfo()

	initializer.Init(
		run,
		initializer.WithMsgf("%s [%s]", config.AppName, config.AppVersion),
		initializer.WithOptionsLogz(logz.WithCaller(false)),
	)
}

func run(ctx context.Context, _ *sync.WaitGroup) error {
	// load config
	if err := config.Load(ctx); err != nil {
		return err //nolint:wrapcheck // no need
	}

	// init telemetry
	collector, err := tell.New(ctx, config.Application.Telemetry)
	if err != nil {
		return fmt.Errorf("failed to init telemetry; %w", err)
	}
	defer collector.Shutdown() //nolint:errcheck // no need

	client, err := wkafka.New(
		ctx,
		config.Application.KafkaConfig,
		wkafka.WithAutoTopicCreation(false),
		wkafka.WithKGOOptions(
			kgo.UnknownTopicRetries(0),
		),
	)
	if err != nil {
		return err //nolint:wrapcheck // no need
	}

	e, err := server.Set(ctx, server.SetConfig{
		Client: client,
	})
	if err != nil {
		return err //nolint:wrapcheck // no need
	}

	// add server shutdown function if context is canceled
	initializer.Shutdown.Add(func() error { return server.Stop(e) }, initializer.WithShutdownName("server")) //nolint:wrapcheck // no need

	// start server
	if err := server.Start(e); err != nil {
		return err //nolint:wrapcheck // no need
	}

	return nil
}
