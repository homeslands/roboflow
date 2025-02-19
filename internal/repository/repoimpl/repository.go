package repoimpl

import (
	"github.com/tuanvumaihuynh/roboflow/internal/db/sqlcpg"
	"github.com/tuanvumaihuynh/roboflow/internal/repository"
)

var _ repository.Repository = (*repoimpl)(nil)

type repoimpl struct {
	qrLocationRepository        *qrLocationRepository
	raybotRepository            *raybotRepository
	raybotCommandRepository     *raybotCommandRepository
	workflowRepository          *workflowRepository
	workflowExecutionRepository *workflowExecutionRepository
	stepExecutionRepository     *stepExecutionRepository
}

//nolint:revive
func NewRepository(queries sqlcpg.Queries) *repoimpl {
	return &repoimpl{
		qrLocationRepository:        newQRLocationRepository(queries),
		raybotRepository:            newRaybotRepository(queries),
		raybotCommandRepository:     newRaybotCommandRepository(queries),
		workflowRepository:          newWorkflowRepository(queries),
		workflowExecutionRepository: newWorkflowExecutionRepository(queries),
		stepExecutionRepository:     newStepExecutionRepository(queries),
	}
}

func (r repoimpl) QRLocation() repository.QRLocationRepository {
	return r.qrLocationRepository
}

func (r repoimpl) Raybot() repository.RaybotRepository {
	return r.raybotRepository
}

func (r repoimpl) RaybotCommand() repository.RaybotCommandRepository {
	return r.raybotCommandRepository
}

func (r repoimpl) Workflow() repository.WorkflowRepository {
	return r.workflowRepository
}

func (r repoimpl) WorkflowExecution() repository.WorkflowExecutionRepository {
	return r.workflowExecutionRepository
}

func (r repoimpl) StepExecution() repository.StepExecutionRepository {
	return r.stepExecutionRepository
}
