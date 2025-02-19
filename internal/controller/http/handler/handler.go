package handler

import (
	"github.com/tuanvumaihuynh/roboflow/internal/controller/http/oas/gen"
	"github.com/tuanvumaihuynh/roboflow/internal/service"
)

var _ gen.StrictServerInterface = (*APIHandler)(nil)

type APIHandler struct {
	*qrLocationHandler
	*raybotHandler
	*raybotCommandHandler
	*workflowHandler
	*workflowExecutionHandler
	*stepExecutionHandler
}

func NewAPIHandler(svc service.Service) *APIHandler {
	return &APIHandler{
		qrLocationHandler:        newQRLocationHandler(svc.QRLocation()),
		raybotHandler:            newRaybotHandler(svc.Raybot()),
		raybotCommandHandler:     newRaybotCommandHandler(svc.RaybotCommand()),
		workflowHandler:          newWorkflowHandler(svc.Workflow()),
		workflowExecutionHandler: newWorkflowExecutionHandler(svc.WorkflowExecution()),
		stepExecutionHandler:     newStepExecutionHandler(svc.StepExecution()),
	}
}
