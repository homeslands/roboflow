package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/tuanvumaihuynh/roboflow/db"
	raybotHttp "github.com/tuanvumaihuynh/roboflow/internal/raybot/port"
	raybotRepo "github.com/tuanvumaihuynh/roboflow/internal/raybot/repository"
	raybotSvc "github.com/tuanvumaihuynh/roboflow/internal/raybot/service"
	"github.com/tuanvumaihuynh/roboflow/pkg/config"
	"github.com/tuanvumaihuynh/roboflow/pkg/middleware"
)

type services struct {
	raybotSvc raybotSvc.RaybotService
}

type CleanupFunc func(ctx context.Context)

func Run(cfg *config.Config, conn *pgxpool.Pool, logger *zap.Logger) CleanupFunc {
	// Setup store and repositories
	store := db.NewStore(conn)
	raybotRepo := raybotRepo.NewRaybotRepository(store)

	// Setup services
	svcs := services{
		raybotSvc: *raybotSvc.NewRaybotService(raybotRepo),
	}

	// Setup server
	srv := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", cfg.Port),
		Handler: setupHandler(chi.NewRouter(), svcs, cfg, logger),
	}

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
	}
}

// setupHandler sets up the HTTP routes and middleware for the server.
func setupHandler(
	r chi.Router,
	svcs services,
	cfg *config.Config,
	logger *zap.Logger,
) chi.Router {
	// Setup handlers
	raybotHandler := raybotHttp.NewRaybotHandler(svcs.raybotSvc, logger)

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
	})

	// Setup health check route
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	return r
}
