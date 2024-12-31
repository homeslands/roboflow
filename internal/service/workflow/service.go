package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/tuanvumaihuynh/roboflow/internal/model"
	"github.com/tuanvumaihuynh/roboflow/internal/service/workflow/event"
	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/pubsub"
	"github.com/tuanvumaihuynh/roboflow/pkg/xerrors"
)

var _ Service = (*service)(nil)

type Service interface {
	Create(ctx context.Context, cmd CreateWorkflowCommand) (model.Workflow, error)
	Update(ctx context.Context, cmd UpdateWorkflowCommand) (model.Workflow, error)
	Delete(ctx context.Context, cmd DeleteWorkflowCommand) error
	Run(ctx context.Context, cmd RunWorkflowCommand) (uuid.UUID, error)

	GetByID(ctx context.Context, q GetWorkflowByIDQuery) (model.Workflow, error)
	List(ctx context.Context, q ListWorkflowQuery) (*paging.List[model.Workflow], error)
}

type service struct {
	workflowRepo          model.WorkflowRepository
	workflowExecutionRepo model.WorkflowExecutionRepository
	eventPublisher        pubsub.Publisher
	log                   *slog.Logger
}

func NewService(workflowRepo model.WorkflowRepository,
	workflowExecutionRepo model.WorkflowExecutionRepository,
	eventPublisher pubsub.Publisher,
	log *slog.Logger,
) *service {
	return &service{
		workflowRepo:          workflowRepo,
		workflowExecutionRepo: workflowExecutionRepo,
		eventPublisher:        eventPublisher,
		log:                   log,
	}
}

func (s service) Create(ctx context.Context, cmd CreateWorkflowCommand) (model.Workflow, error) {
	if err := cmd.Validate(); err != nil {
		return model.Workflow{}, err
	}

	now := time.Now()
	modelWorkflow := model.Workflow{
		ID:          uuid.New(),
		Name:        cmd.Name,
		Description: cmd.Description,
		Definition:  &cmd.Definition,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// TODO: validate definition
	// - Check node, edge
	// - Check node env, workflow env

	if err := s.workflowRepo.Create(ctx, modelWorkflow); err != nil {
		return model.Workflow{}, err
	}

	return modelWorkflow, nil
}

func (s service) Update(ctx context.Context, cmd UpdateWorkflowCommand) (model.Workflow, error) {
	if err := cmd.Validate(); err != nil {
		return model.Workflow{}, err
	}

	wf, err := s.workflowRepo.Update(ctx, model.Workflow{
		ID:          cmd.ID,
		Name:        cmd.Name,
		Description: cmd.Description,
		Definition:  &cmd.Definition,
		UpdatedAt:   time.Now(),
	})
	if err != nil {
		return model.Workflow{}, err
	}

	return wf, nil
}

func (s service) Delete(ctx context.Context, cmd DeleteWorkflowCommand) error {
	if err := cmd.Validate(); err != nil {
		return err
	}

	return s.workflowRepo.Delete(ctx, cmd.ID)
}

func (s service) Run(ctx context.Context, cmd RunWorkflowCommand) (uuid.UUID, error) {
	if err := cmd.Validate(); err != nil {
		return uuid.Nil, err
	}

	wf, err := s.workflowRepo.Get(ctx, cmd.ID)
	if err != nil {
		return uuid.Nil, xerrors.ThrowNotFound(err, "workflow not found")
	}

	// Workflow Execution ID
	wfeID := uuid.New()

	// Create steps
	steps, err := s.createSteps(ctx, wf.Definition.Nodes, cmd.Env, wfeID)
	if err != nil {
		return uuid.Nil, err
	}

	// Create workflow execution
	wfe := model.WorkflowExecution{
		ID:         wfeID,
		WorkflowID: wf.ID,
		Status:     model.WorkflowExecutionStatusPending,
		Env:        cmd.Env,
		Definition: *wf.Definition,
		CreatedAt:  time.Now(),
		Steps:      &steps,
	}
	err = s.workflowExecutionRepo.Create(ctx, wfe)
	if err != nil {
		return uuid.Nil, err
	}

	// Publish event
	s.log.DebugContext(ctx, "Publish workflow run event",
		slog.String("workflow_execution_id", wfeID.String()))

	msgJSON, err := json.Marshal(wfe)
	if err != nil {
		s.log.ErrorContext(ctx, "Failed to create workflow run event",
			slog.String("workflow_execution_id", wfeID.String()),
			slog.Any("error", err))
	}

	if err := s.eventPublisher.Publish(event.TopicWorkflowExecutionRun, msgJSON); err != nil {
		s.log.ErrorContext(ctx, "Failed to publish workflow run event",
			slog.String("topic", event.TopicWorkflowExecutionRun),
			slog.Any("error", err))
	}

	return wfeID, nil
}

func (s service) GetByID(ctx context.Context, q GetWorkflowByIDQuery) (model.Workflow, error) {
	if err := q.Validate(); err != nil {
		return model.Workflow{}, err
	}

	return s.workflowRepo.Get(ctx, q.ID)
}

func (s service) List(ctx context.Context, q ListWorkflowQuery) (*paging.List[model.Workflow], error) {
	if err := q.Validate(); err != nil {
		return nil, err
	}

	return s.workflowRepo.List(ctx, q.PagingParams, q.Sorts)
}

func (s service) createSteps(ctx context.Context, nodes []model.WorkflowNode, env map[string]string, wfeID uuid.UUID) ([]model.Step, error) {
	steps := make([]model.Step, 0, len(nodes))
	for _, node := range nodes {
		stepEnv := make(map[string]string)
		for sysKey, field := range node.Definition.Fields {
			if field.UseEnv {
				if field.Key == nil {
					return nil, xerrors.ThrowInternal(nil,
						fmt.Sprintf("env key field is required when use env in node id %s", node.ID),
					)
				}
				val, exists := env[*field.Key]
				if !exists {
					return nil, xerrors.ThrowInvalidArgument(nil,
						fmt.Sprintf("env key %s not found in environment variable", *field.Key),
						xerrors.WithCode("missing_env"),
					)
				}
				field.Value = &val
				node.Definition.Fields[sysKey] = field
				stepEnv[sysKey] = val
			}
		}

		step := model.Step{
			ID:                  uuid.New(),
			WorkflowExecutionID: wfeID,
			Env:                 stepEnv,
			Node:                node,
			Status:              model.WorkflowExecutionStepStatusPending,
		}
		steps = append(steps, step)

		s.log.DebugContext(ctx, "Create step",
			slog.String("step_id", step.ID.String()),
			slog.Any("env", stepEnv))
	}
	return steps, nil
}
