package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	pubsubMiddleware "github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/tuanvumaihuynh/roboflow/db"
	commandEvent "github.com/tuanvumaihuynh/roboflow/internal/command/event"
	commandHttp "github.com/tuanvumaihuynh/roboflow/internal/command/port"
	commandRepo "github.com/tuanvumaihuynh/roboflow/internal/command/repository"
	commandSvc "github.com/tuanvumaihuynh/roboflow/internal/command/service"
	qrLocationHttp "github.com/tuanvumaihuynh/roboflow/internal/location/port"
	qrLocationRepo "github.com/tuanvumaihuynh/roboflow/internal/location/repository"
	qrLocationSvc "github.com/tuanvumaihuynh/roboflow/internal/location/service"
	raybotHttp "github.com/tuanvumaihuynh/roboflow/internal/raybot/port"
	raybotRepo "github.com/tuanvumaihuynh/roboflow/internal/raybot/repository"
	raybotSvc "github.com/tuanvumaihuynh/roboflow/internal/raybot/service"
	"github.com/tuanvumaihuynh/roboflow/internal/ws/hub"
	"github.com/tuanvumaihuynh/roboflow/pkg/config"
	"github.com/tuanvumaihuynh/roboflow/pkg/middleware"
	"github.com/tuanvumaihuynh/roboflow/pkg/pubsub"
)

type services struct {
	raybotSvc     raybotSvc.RaybotService
	qrLocationSvc qrLocationSvc.QrLocationService
	commandSvc    commandSvc.CommandService
}

type CleanupFunc func(ctx context.Context)

func Run(cfg *config.Config, conn *pgxpool.Pool, logger *zap.Logger) CleanupFunc {
	// Setup pubsub
	pubsubLogger := watermill.NewStdLogger(false, true)
	gChan := gochannel.NewGoChannel(gochannel.Config{}, pubsubLogger)
	publisher := pubsub.NewPublisher(gChan)
	pubsubRouter, _ := message.NewRouter(message.RouterConfig{}, pubsubLogger)
	pubsubRouter.AddMiddleware(
		pubsubMiddleware.Retry{
			MaxRetries:      1,
			InitialInterval: time.Millisecond * 100,
			Logger:          pubsubLogger,
		}.Middleware,
		pubsubMiddleware.Recoverer,
	)

	// Setup store and repositories
	store := db.NewStore(conn)
	raybotRepo := raybotRepo.NewRaybotRepository(store)
	qrLocationRepo := qrLocationRepo.NewQrRepository(store)
	commandRepo := commandRepo.NewMemoryCommandRepository()

	// Setup services
	svcs := services{
		raybotSvc:     *raybotSvc.NewRaybotService(raybotRepo),
		qrLocationSvc: *qrLocationSvc.NewQrLocationService(qrLocationRepo),
		commandSvc:    *commandSvc.NewCommandService(commandRepo, publisher),
	}

	// Setup websocket hub
	hub := hub.NewHub(hub.HubConfig{Logger: logger})

	// Setup event handler
	pubsubRouter = setupEventHandler(pubsubRouter, gChan, hub)

	// Setup server
	srv := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", cfg.Port),
		Handler: setupHandler(chi.NewRouter(), hub, svcs, cfg, logger),
	}

	// Run pubsub router
	go func() {
		if err := pubsubRouter.Run(context.Background()); err != nil {
			logger.Error("Error running pubsub router", zap.Error(err))
		}
	}()

	// Run api server
	go func() {
		logger.Sugar().Infof("Server started at %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Error starting server", zap.Error(err))
		}
	}()

	return func(ctx context.Context) {
		logger.Info("Shutting down server")
		// Shutdown server
		if err := srv.Shutdown(ctx); err != nil {
			logger.Error("Error shutting down server", zap.Error(err))
		}
		// Shutdown pubsub router
		if err := pubsubRouter.Close(); err != nil {
			logger.Error("Error closing pubsub router", zap.Error(err))
		}
	}
}

// setupHandler sets up the HTTP routes and middleware for the server.
func setupHandler(
	r chi.Router,
	hub *hub.Hub,
	svcs services,
	cfg *config.Config,
	logger *zap.Logger,
) chi.Router {
	// Setup handlers
	raybotHandler := raybotHttp.NewRaybotHandler(svcs.raybotSvc, logger)
	qrLocationHandler := qrLocationHttp.NewQrLocationHandler(svcs.qrLocationSvc, logger)
	commandHandler := commandHttp.NewCommandHandler(svcs.commandSvc, logger)

	// Setup middleware
	r.Use(middleware.Logging(logger))
	r.Use(middleware.Cors())
	r.Use(func(n http.Handler) http.Handler {
		return middleware.Recoverer(n, cfg.IsProd(), logger)
	})

	// Setup API Routes
	r.Route("/api", func(r chi.Router) {
		r.Route("/raybots", func(r chi.Router) {
			r.Get("/{id}", raybotHandler.HandleGetRaybot)
			r.Get("/", raybotHandler.HandleListRaybots)
			r.Post("/", raybotHandler.HandleCreateRaybot)
			r.Delete("/{id}", raybotHandler.HandleDeleteRaybot)
		})
		r.Route("/qr-locations", func(r chi.Router) {
			r.Get("/{id}", qrLocationHandler.HandleGetQRLocation)
			r.Get("/", qrLocationHandler.HandleListQRLocations)
			r.Post("/", qrLocationHandler.HandleCreateQRLocation)
			r.Put("/{id}", qrLocationHandler.HandleUpdateQRLocation)
			r.Delete("/{id}", qrLocationHandler.HandleDeleteQRLocation)
		})
		r.Route("/commands", func(r chi.Router) {
			r.Get("/{id}", commandHandler.HandleGetCommand)
			r.Get("/", commandHandler.HandleListCommands)
			r.Post("/", commandHandler.HandleCreateCommand)
			r.Delete("/{id}", commandHandler.HandleDeleteCommand)
		})
	})

	// Setup websocket routes
	r.Get("/ws-raybot", func(w http.ResponseWriter, r *http.Request) {
		hub.HandleRaybotClientConnect(w, r, svcs.raybotSvc, svcs.commandSvc)
	})

	// Setup health check route
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	return r
}

func setupEventHandler(
	r *message.Router,
	gChan *gochannel.GoChannel,
	hub *hub.Hub,
) *message.Router {
	r.AddNoPublisherHandler(
		"command-to-ws-hub",
		commandEvent.TopicCommandCreated,
		gChan,
		func(msg *message.Message) error {
			return hub.HandleCommandCreated(msg)
		},
	)

	return r
}
