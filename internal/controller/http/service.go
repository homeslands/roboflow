package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"

	httphandler "github.com/tuanvumaihuynh/roboflow/internal/controller/http/handler"
	"github.com/tuanvumaihuynh/roboflow/internal/controller/http/oas/gen"
	"github.com/tuanvumaihuynh/roboflow/internal/controller/http/swagger"
	"github.com/tuanvumaihuynh/roboflow/internal/service"
	"github.com/tuanvumaihuynh/roboflow/pkg/config"
	"github.com/tuanvumaihuynh/roboflow/pkg/middleware"
)

//nolint:revive
type HTTPService struct {
	config  config.HTTPServerConfig
	service service.Service
	log     *slog.Logger
}

func NewHTTPService(
	config config.HTTPServerConfig,
	service service.Service,
	log *slog.Logger,
) *HTTPService {
	return &HTTPService{
		config:  config,
		service: service,
		log:     log.With(slog.String("service", "http_service")),
	}
}

type CleanupFunc func(ctx context.Context) error

func (s HTTPService) Run() (CleanupFunc, error) {
	r := chi.NewRouter()

	s.registerMiddlewares(r)

	if s.config.EnableSwagger {
		s.registerSwaggerHandler(r)
	}

	s.registerAPIHandler(r)

	return s.RunWithServer(r)
}

func (s HTTPService) RunWithServer(r chi.Router) (CleanupFunc, error) {
	if err := chi.Walk(r, func(method, route string, _ http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		s.log.Debug(fmt.Sprintf("%7s '%s' has %d middlewares", method, route, len(middlewares)))
		return nil
	}); err != nil {
		return nil, fmt.Errorf("error walking HTTP router: %w", err)
	}

	srv := &http.Server{
		Addr:              net.JoinHostPort("0.0.0.0", fmt.Sprintf("%d", s.config.Port)),
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		s.log.Info(fmt.Sprintf("starting HTTP server at %s", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.log.Error("error starting HTTP server", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	cleanup := func(ctx context.Context) error {
		if err := srv.Shutdown(ctx); err != nil {
			s.log.Error("error shutting down HTTP server", slog.String("error", err.Error()))
			return err
		}

		return nil
	}

	return cleanup, nil
}

func (s HTTPService) registerAPIHandler(r chi.Router) {
	apiHandler := httphandler.NewAPIHandler(s.service)
	strictAPIHandler := gen.NewStrictHandlerWithOptions(
		apiHandler,
		[]gen.StrictMiddlewareFunc{},
		gen.StrictHTTPServerOptions{
			RequestErrorHandlerFunc:  s.handleRequestError,
			ResponseErrorHandlerFunc: s.handleResponseError,
		},
	)

	gen.HandlerFromMuxWithBaseURL(strictAPIHandler, r, "/api/v1")
}

func (s HTTPService) registerMiddlewares(r chi.Router) {
	r.Use(middleware.RequestID)
	r.Use(middleware.Logging)
	r.Use(middleware.Cors())
	r.Use(middleware.Recoverer)
}

func (s HTTPService) registerSwaggerHandler(r chi.Router) {
	swagger.Register(r, "/docs/openapi.yml")
}
