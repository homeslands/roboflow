package application

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	pubsubMiddleware "github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/tuanvumaihuynh/roboflow/db"
	"github.com/tuanvumaihuynh/roboflow/docs/openapi"
	"github.com/tuanvumaihuynh/roboflow/internal/api"
	"github.com/tuanvumaihuynh/roboflow/internal/repository"
	qrLocationSvc "github.com/tuanvumaihuynh/roboflow/internal/service/qr_location"
	raybotSvc "github.com/tuanvumaihuynh/roboflow/internal/service/raybot"
	raybotCommandSvc "github.com/tuanvumaihuynh/roboflow/internal/service/raybot_command"
	"github.com/tuanvumaihuynh/roboflow/internal/service/raybot_command/event"
	stepSvc "github.com/tuanvumaihuynh/roboflow/internal/service/step"
	workflowSvc "github.com/tuanvumaihuynh/roboflow/internal/service/workflow"
	workflowexecution "github.com/tuanvumaihuynh/roboflow/internal/service/workflow_execution"
	raybotclient "github.com/tuanvumaihuynh/roboflow/internal/ws/raybot_client"
	"github.com/tuanvumaihuynh/roboflow/pkg/config"
	"github.com/tuanvumaihuynh/roboflow/pkg/middleware"
	"github.com/tuanvumaihuynh/roboflow/pkg/pubsub"
)

type application struct {
	cfg            config.Config
	log            *slog.Logger
	eventPublisher pubsub.Publisher

	raybotSvc            raybotSvc.Service
	raybotCommandSvc     raybotCommandSvc.Service
	qrLocationSvc        qrLocationSvc.Service
	workflowSvc          workflowSvc.Service
	workflowExecutionSvc workflowexecution.Service
	stepSvc              stepSvc.Service

	wsRaybot *raybotclient.WebSocket
}

type CleanupFunc func(ctx context.Context)

func Run(cfg config.Config, conn *pgxpool.Pool, logger *slog.Logger) CleanupFunc {
	// Setup pubsub
	pubsubLogger := watermill.NewSlogLogger(logger.With(slog.String("service", "pubsub")))
	gChan := gochannel.NewGoChannel(
		gochannel.Config{
			BlockPublishUntilSubscriberAck: false,
		},
		pubsubLogger,
	)
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
	raybotRepo := repository.NewRaybotRepository(store)
	raybotCommandRepo := repository.NewRaybotCommandRepository(store)
	qrLocationRepo := repository.NewQRLocationRepository(store)
	workflowRepo := repository.NewWorkflowRepository(store)
	workflowExecutionRepo := repository.NewWorkflowExecutionRepository(store)
	stepRepo := repository.NewStepRepository(store)

	// Setup services
	raybotSvc := raybotSvc.NewService(raybotRepo)
	raybotCommandSvc := raybotCommandSvc.NewService(raybotCommandRepo, raybotRepo, qrLocationRepo, publisher, logger)
	qrLocationSvc := qrLocationSvc.NewService(qrLocationRepo)
	workflowSvc := workflowSvc.NewService(workflowRepo, workflowExecutionRepo, publisher, logger)
	workflowExecutionSvc := workflowexecution.NewService(workflowExecutionRepo)
	stepSvc := stepSvc.NewService(stepRepo)

	// Setup websocket
	raybotWs := raybotclient.NewWebSocket(raybotRepo, raybotCommandSvc, logger)

	// Setup application
	app := application{
		cfg:            cfg,
		log:            logger,
		eventPublisher: publisher,

		raybotSvc:            raybotSvc,
		raybotCommandSvc:     raybotCommandSvc,
		qrLocationSvc:        qrLocationSvc,
		workflowSvc:          workflowSvc,
		workflowExecutionSvc: workflowExecutionSvc,
		stepSvc:              stepSvc,

		wsRaybot: raybotWs,
	}

	// Setup server
	srv := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", cfg.Port),
		Handler: setupAPIHandler(app),
	}

	// Run pubsub router
	setupEventHandler(pubsubRouter, gChan, app)
	go func() {
		if err := pubsubRouter.Run(context.Background()); err != nil {
			logger.Error("Error running pubsub router", slog.Any("error", err))
		}
	}()

	// Run api server
	go func() {
		logger.Info("Starting server", slog.Int("port", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Error starting server", slog.Any("error", err))
		}
	}()

	return func(ctx context.Context) {
		logger.Info("Shutting down server")
		// Shutdown server
		if err := srv.Shutdown(ctx); err != nil {
			logger.Error("Error shutting down server", slog.Any("error", err))
		}
		// Shutdown pubsub router
		if err := pubsubRouter.Close(); err != nil {
			logger.Error("Error closing pubsub router", slog.Any("error", err))
		}
	}
}

func setupAPIHandler(
	app application,
) chi.Router {
	r := chi.NewRouter()

	// Setup middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.Logging(app.log))
	r.Use(middleware.Cors())
	r.Use(func(n http.Handler) http.Handler {
		return middleware.Recoverer(n, app.cfg.IsProd(), app.log)
	})

	// Setup HTTP server
	h := api.NewHTTPServer(
		app.log,
		app.raybotSvc, app.qrLocationSvc, app.raybotCommandSvc,
		app.workflowSvc, app.workflowExecutionSvc, app.stepSvc,
	)
	w := api.ServerInterfaceWrapper{
		Handler:          h,
		ErrorHandlerFunc: h.HandleError,
	}
	r.Route("/api/v1", func(r chi.Router) {
		// Raybot
		r.Get("/raybots", w.ListRaybots)
		r.Post("/raybots", w.CreateRaybot)
		r.Get("/raybots/{raybotId}", w.GetRaybotById)
		r.Delete("/raybots/{raybotId}", w.DeleteRaybotById)

		// Raybot Command
		r.Get("/raybots/{raybotId}/commands", w.ListRaybotCommands)
		r.Post("/raybots/{raybotId}/commands", w.CreateRaybotCommand)
		r.Get("/raybot-commands/{raybotCommandId}", w.GetRaybotCommandById)

		// QR Location
		r.Get("/qr-locations", w.ListQRLocations)
		r.Post("/qr-locations", w.CreateQRLocation)
		r.Get("/qr-locations/{qrLocationId}", w.GetQRLocationById)
		r.Put("/qr-locations/{qrLocationId}", w.UpdateQRLocationById)
		r.Delete("/qr-locations/{qrLocationId}", w.DeleteQRLocationById)

		// Workflow
		r.Get("/workflows", w.ListWorkflows)
		r.Post("/workflows", w.CreateWorkflow)
		r.Get("/workflows/{workflowId}", w.GetWorkflowById)
		r.Put("/workflows/{workflowId}", w.UpdateWorkflowById)
		r.Delete("/workflows/{workflowId}", w.DeleteWorkflowById)
		r.Post("/workflows/{workflowId}/run", w.RunWorkflowById)

		// Workflow Execution
		r.Get("/workflows/{workflowId}/executions", w.ListWorkflowExecutionsByWorkflowId)
		r.Get("/workflow-executions/{workflowExecutionId}", w.GetWorkflowExecutionById)
		r.Get("/workflow-executions/{workflowExecutionId}/status", w.GetWorkflowExecutionStatusById)

		// Step
		r.Get("/workflow-executions/{workflowExecutionId}/steps", w.ListStepsByWorkflowExecutionId)
		r.Get("/steps/{stepId}", w.GetStepById)
	})
	r.Get("/health", h.HandleHealth)
	r.NotFound(h.HandleNotFound)
	r.MethodNotAllowed(h.HandleMethodNotAllowed)
	setupOpenapiDocs(r)

	// api.HandlerFromMuxWithBaseURL(httpServer, r, "/api/v1")

	// Setup websocket
	app.wsRaybot.RegisterHandlers(r)

	return r
}

func setupOpenapiDocs(r chi.Router) {
	spec, err := openapi.OpenapiSpec.ReadFile("build/openapi.yml")
	if err != nil {
		panic(err)
	}
	specPath := "/docs/openapi.json"
	template := openapi.GetTemplate(specPath)

	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(template))
	})
	r.Get(specPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yaml")
		w.Write(spec)
	})
}

func setupEventHandler(
	r *message.Router,
	pubsub *gochannel.GoChannel,
	app application,
) {
	r.AddNoPublisherHandler(
		"raybot_command_to_ws_raybot_client",
		event.TopicRaybotCommandCreated,
		pubsub,
		func(msg *message.Message) error {
			return app.wsRaybot.HandleRaybotCommandCreated(msg.Payload)
		},
	)
}
