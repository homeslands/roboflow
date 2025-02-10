package serviceimpl

import (
	"log/slog"

	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/tuanvumaihuynh/roboflow/internal/db/sqldb"
	"github.com/tuanvumaihuynh/roboflow/internal/repository"
	"github.com/tuanvumaihuynh/roboflow/internal/service"
	"github.com/tuanvumaihuynh/roboflow/pkg/validator"
)

var _ service.Service = (*serviceimpl)(nil)

type serviceimpl struct {
	qrLocationService        *qrLocationService
	raybotService            *raybotService
	raybotCommandService     *raybotCommandService
	workflowService          *workflowService
	workflowExecutionService *workflowExecutionService
	stepExecutionService     *stepExecutionService
}

//nolint:revive
func NewService(
	repository repository.Repository,
	sqlDBProvider sqldb.Provider,
	publisher message.Publisher,
	validator validator.Validator,
	log *slog.Logger,
) *serviceimpl {
	qrLocationSvc := newQRLocationService(repository.QRLocation(), sqlDBProvider, validator)
	raybotSvc := newRaybotService(repository.Raybot(), sqlDBProvider, validator)
	raybotCommandSvc := newRaybotCommandService(repository.RaybotCommand(), sqlDBProvider, publisher, validator)
	workflowSvc := newWorkflowService(repository.Workflow(), repository.WorkflowExecution(),
		repository.StepExecution(), sqlDBProvider, publisher, validator)
	workflowExecutionSvc := newWorkflowExecutionService(repository.WorkflowExecution(),
		repository.StepExecution(), sqlDBProvider, validator)
	stepExecutionSvc := newStepExecutionService(repository.StepExecution(), sqlDBProvider, validator)

	return &serviceimpl{
		qrLocationService:        qrLocationSvc,
		raybotService:            raybotSvc,
		raybotCommandService:     raybotCommandSvc,
		workflowService:          workflowSvc,
		workflowExecutionService: workflowExecutionSvc,
		stepExecutionService:     stepExecutionSvc,
	}
}

func (s *serviceimpl) QRLocation() service.QRLocationService {
	return s.qrLocationService
}

func (s *serviceimpl) Raybot() service.RaybotService {
	return s.raybotService
}

func (s *serviceimpl) RaybotCommand() service.RaybotCommandService {
	return s.raybotCommandService
}

func (s *serviceimpl) Workflow() service.WorkflowService {
	return s.workflowService
}

func (s *serviceimpl) WorkflowExecution() service.WorkflowExecutionService {
	return s.workflowExecutionService
}

func (s *serviceimpl) StepExecution() service.StepExecutionService {
	return s.stepExecutionService
}
