package application

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/tuanvumaihuynh/roboflow/internal/db/sqlcpg"
	"github.com/tuanvumaihuynh/roboflow/internal/db/sqldb"
	"github.com/tuanvumaihuynh/roboflow/internal/pubsub"
	"github.com/tuanvumaihuynh/roboflow/internal/repository/repoimpl"
	"github.com/tuanvumaihuynh/roboflow/internal/service"
	"github.com/tuanvumaihuynh/roboflow/internal/service/serviceimpl"
	"github.com/tuanvumaihuynh/roboflow/pkg/config"
	mylog "github.com/tuanvumaihuynh/roboflow/pkg/log"
	"github.com/tuanvumaihuynh/roboflow/pkg/pgxslog"
	"github.com/tuanvumaihuynh/roboflow/pkg/validator"
)

type Application struct {
	Config *config.Config

	Service service.Service

	Publisher  message.Publisher
	Subscriber message.Subscriber

	Log *slog.Logger

	context context.Context
}

func (a *Application) Context() context.Context {
	return a.context
}

type CleanupFunc func() error

func New() (*Application, CleanupFunc) {
	// Load config from environment variables
	conf, err := config.Load()
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	// Set UTC timezone
	time.Local = time.UTC

	// Context
	ctx := context.Background()

	// Init logger
	log := mylog.NewLogger(mylog.Config{
		Level:     conf.Log.Level,
		Format:    conf.Log.Format,
		AddSource: conf.Log.AddSource,
	})
	slog.SetDefault(log)

	// Setup postgres
	pgConf, err := pgxpool.ParseConfig(conf.Postgres.ConnectionString())
	if err != nil {
		log.Error("error parsing postgres connection string", slog.Any("error", err))
		os.Exit(1)
	}

	if conf.Postgres.EnableLog {
		pgConf.ConnConfig.Tracer = pgxslog.NewTracer(log)
	}
	// TODO: Setup otelpgx

	if conf.Postgres.MaxConns > 0 {
		pgConf.MaxConns = conf.Postgres.MaxConns
	}
	if conf.Postgres.MinConns > 0 {
		pgConf.MinConns = conf.Postgres.MinConns
	}
	pgConf.MaxConnLifetime = 15 * 60 * time.Second

	pgPool, err := pgxpool.NewWithConfig(ctx, pgConf)
	if err != nil {
		log.Error("error creating postgres pool", slog.Any("error", err))
		os.Exit(1)
	}

	// Setup Nats pubsub
	publisher, subscriber, err := pubsub.NewNatsPubSub(conf.Nats, log)
	if err != nil {
		log.Error("error creating nats pubsub", slog.Any("error", err))
		os.Exit(1)
	}

	// Setup repository
	sqlDBProvider := sqldb.NewProvider(pgPool)
	queries := sqlcpg.New()
	repo := repoimpl.NewRepository(*queries)

	// Setup service
	validator := validator.NewValidator()
	svc := serviceimpl.NewService(repo, sqlDBProvider, publisher, validator, log)

	// Setup application
	app := &Application{
		Config:     conf,
		Service:    svc,
		Publisher:  publisher,
		Subscriber: subscriber,
		Log:        log,
		context:    ctx,
	}

	// Cleanup function
	cleanup := func() error {
		if err := publisher.Close(); err != nil {
			return fmt.Errorf("error closing publisher: %w", err)
		}
		if err := subscriber.Close(); err != nil {
			return fmt.Errorf("error closing subscriber: %w", err)
		}
		pgPool.Close()
		return nil
	}

	return app, cleanup
}
