package api

import (
	"fmt"

	"github.com/tuanvumaihuynh/roboflow/internal/application"
	"github.com/tuanvumaihuynh/roboflow/internal/controller/http"
)

func Start(app *application.Application, interruptChan <-chan any) error {
	httpSvc := http.NewHTTPService(app.Config.HTTPServer, app.Service, app.Log)

	cleanup, err := httpSvc.Run()
	if err != nil {
		return fmt.Errorf("error running http server: %w", err)
	}

	<-interruptChan

	app.Log.Debug("http server shutting down")

	if err := cleanup(app.Context()); err != nil {
		return fmt.Errorf("error cleaning up http server: %w", err)
	}

	app.Log.Debug("http server shutdown complete")

	return nil
}
