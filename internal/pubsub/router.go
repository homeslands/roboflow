package pubsub

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
)

// NewRouter creates a new pubsub router
func NewRouter(log *slog.Logger) (*message.Router, error) {
	wLog := watermill.NewSlogLogger(log)

	router, err := message.NewRouter(message.RouterConfig{
		CloseTimeout: time.Second * 30,
	}, wLog)
	if err != nil {
		return nil, fmt.Errorf("create router: %w", err)
	}

	// Add middlewares
	router.AddMiddleware(
		middleware.Recoverer,
		middleware.CorrelationID,
		middleware.Retry{
			MaxRetries:      3,
			InitialInterval: time.Second,
			Logger:          wLog,
		}.Middleware,
	)

	return router, nil
}
