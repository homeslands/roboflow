package api

import (
	"log/slog"

	qrlocation "github.com/tuanvumaihuynh/roboflow/internal/service/qr_location"
	"github.com/tuanvumaihuynh/roboflow/internal/service/raybot"
	raybotcommand "github.com/tuanvumaihuynh/roboflow/internal/service/raybot_command"
	"github.com/tuanvumaihuynh/roboflow/internal/service/step"
	"github.com/tuanvumaihuynh/roboflow/internal/service/workflow"
	workflowexecution "github.com/tuanvumaihuynh/roboflow/internal/service/workflow_execution"
)

var _ ServerInterface = (*HTTPServer)(nil)

type HTTPServer struct {
	log *slog.Logger

	raybotSvc            raybot.Service
	qrLocationSvc        qrlocation.Service
	raybotCommandSvc     raybotcommand.Service
	workflowSvc          workflow.Service
	workflowExecutionSvc workflowexecution.Service
	stepSvc              step.Service
}

func NewHTTPServer(
	log *slog.Logger,
	raybotSvc raybot.Service,
	qrLocationSvc qrlocation.Service,
	raybotCommandSvc raybotcommand.Service,
	workflowSvc workflow.Service,
	workflowExecutionSvc workflowexecution.Service,
	stepSvc step.Service,
) *HTTPServer {
	return &HTTPServer{
		log:                  log.With(slog.String("service", "HTTPServer")),
		raybotSvc:            raybotSvc,
		qrLocationSvc:        qrLocationSvc,
		raybotCommandSvc:     raybotCommandSvc,
		workflowSvc:          workflowSvc,
		workflowExecutionSvc: workflowExecutionSvc,
		stepSvc:              stepSvc,
	}
}
