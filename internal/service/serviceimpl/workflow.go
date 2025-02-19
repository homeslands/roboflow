package serviceimpl

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"

	"github.com/tuanvumaihuynh/roboflow/internal/db/sqldb"
	stepexecution "github.com/tuanvumaihuynh/roboflow/internal/model/step_execution"
	"github.com/tuanvumaihuynh/roboflow/internal/model/workflow"
	"github.com/tuanvumaihuynh/roboflow/internal/model/workflow/node"
	workflowexecution "github.com/tuanvumaihuynh/roboflow/internal/model/workflow_execution"
	"github.com/tuanvumaihuynh/roboflow/internal/pubsub"
	"github.com/tuanvumaihuynh/roboflow/internal/repository"
	"github.com/tuanvumaihuynh/roboflow/internal/service"
	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/validator"
	"github.com/tuanvumaihuynh/roboflow/pkg/xerror"
)

var (
	_ service.WorkflowService = (*workflowService)(nil)
)

type workflowService struct {
	workflowRepo          repository.WorkflowRepository
	workflowExecutionRepo repository.WorkflowExecutionRepository
	stepExecutionRepo     repository.StepExecutionRepository
	sqlDBProvider         sqldb.Provider
	publisher             message.Publisher
	validator             validator.Validator
}

func newWorkflowService(
	workflowRepo repository.WorkflowRepository,
	workflowExecutionRepo repository.WorkflowExecutionRepository,
	stepExecutionRepo repository.StepExecutionRepository,
	sqlDBProvider sqldb.Provider,
	publisher message.Publisher,
	validator validator.Validator,
) *workflowService {
	return &workflowService{
		workflowRepo:          workflowRepo,
		workflowExecutionRepo: workflowExecutionRepo,
		stepExecutionRepo:     stepExecutionRepo,
		sqlDBProvider:         sqlDBProvider,
		publisher:             publisher,
		validator:             validator,
	}
}

func (s workflowService) GetWorkflow(ctx context.Context, params service.GetWorkflowParams) (workflow.Workflow, error) {
	if err := s.validator.Validate(params); err != nil {
		return workflow.Workflow{}, fmt.Errorf("validate params: %w", err)
	}

	wf, err := s.workflowRepo.GetWorkflow(ctx, s.sqlDBProvider.DB(), params.ID)
	if err != nil {
		return workflow.Workflow{}, fmt.Errorf("repo get workflow: %w", err)
	}

	return wf, nil
}

func (s workflowService) ListWorkflows(ctx context.Context, params service.ListWorkflowsParams) (paging.List[workflow.Workflow], error) {
	if err := s.validator.Validate(params); err != nil {
		return paging.List[workflow.Workflow]{}, fmt.Errorf("validate params: %w", err)
	}

	wfs, err := s.workflowRepo.ListWorkflows(ctx, s.sqlDBProvider.DB(), params.PagingParams, params.Sorts)
	if err != nil {
		return paging.List[workflow.Workflow]{}, fmt.Errorf("repo list workflows: %w", err)
	}

	return wfs, nil
}

func (s workflowService) CreateWorkflow(ctx context.Context, params service.CreateWorkflowParams) (workflow.Workflow, error) {
	if err := s.validator.Validate(params); err != nil {
		return workflow.Workflow{}, fmt.Errorf("validate params: %w", err)
	}

	isValid := true
	if err := params.Data.Validate(); err != nil {
		isValid = false
	}

	wf := workflow.NewWorkflow(params.Name, params.Description, isValid, params.Data)
	err := s.workflowRepo.CreateWorkflow(ctx, s.sqlDBProvider.DB(), wf)
	if err != nil {
		return workflow.Workflow{}, fmt.Errorf("repo create workflow: %w", err)
	}

	return wf, nil
}

func (s workflowService) UpdateWorkflow(ctx context.Context, params service.UpdateWorkflowParams) (workflow.Workflow, error) {
	if err := s.validator.Validate(params); err != nil {
		return workflow.Workflow{}, fmt.Errorf("validate params: %w", err)
	}

	isValid := true
	if params.SetData {
		if err := params.Data.Validate(); err != nil {
			isValid = false
		}
	}

	wf, err := s.workflowRepo.UpdateWorkflow(ctx, s.sqlDBProvider.DB(), repository.UpdateWorkflowParams{
		ID:             params.ID,
		Name:           params.Name,
		SetName:        params.SetName,
		Description:    params.Description,
		SetDescription: params.SetDescription,
		IsDraft:        params.IsDraft,
		SetIsDraft:     params.SetIsDraft,
		IsValid:        isValid,
		SetIsValid:     true,
		Data:           params.Data,
		SetData:        params.SetData,
	})
	if err != nil {
		return workflow.Workflow{}, fmt.Errorf("repo update workflow: %w", err)
	}

	return wf, nil
}

func (s workflowService) DeleteWorkflow(ctx context.Context, params service.DeleteWorkflowParams) error {
	if err := s.validator.Validate(params); err != nil {
		return fmt.Errorf("validate params: %w", err)
	}

	if err := s.workflowRepo.DeleteWorkflow(ctx, s.sqlDBProvider.DB(), params.ID); err != nil {
		return fmt.Errorf("repo delete workflow: %w", err)
	}

	return nil
}

func (s workflowService) RunWorkflow(ctx context.Context, params service.RunWorkflowParams) (string, error) {
	if err := s.validator.Validate(params); err != nil {
		return "", fmt.Errorf("validate params: %w", err)
	}

	wf, err := s.workflowRepo.GetWorkflow(ctx, s.sqlDBProvider.DB(), params.ID)
	if err != nil {
		return "", fmt.Errorf("repo get workflow: %w", err)
	}

	if wf.IsDraft {
		return "", xerror.ValidationFailed(nil, "Workflow is in draft mode")
	}
	if !wf.IsValid {
		return "", xerror.ValidationFailed(nil, "Workflow is not valid")
	}

	// Get trigger node
	var triggerNode *node.Node
	for _, n := range wf.Data.Nodes {
		if n.Type == node.TypeTrigger {
			triggerNode = &n
			break
		}
	}

	if triggerNode == nil {
		return "", xerror.ValidationFailed(nil, "Trigger node not found")
	}
	triggerData, err := triggerNode.Data.AsTriggerData()
	if err != nil {
		return "", xerror.ValidationFailed(nil, "Trigger node data is not a valid trigger data")
	}

	switch triggerData.TriggerType {
	case node.TriggerTypeOnDemand:
		// Validate runtime variables
		onDemandTriggerData, err := triggerData.AsOnDemandTriggerData()
		if err != nil {
			return "", xerror.ValidationFailed(nil, "Invalid runtime variables")
		}
		for _, v := range onDemandTriggerData.RuntimeVariables {
			if v.Required && params.RuntimeVariables[v.Key] == nil {
				return "", xerror.ValidationFailed(nil, fmt.Sprintf("Runtime variable %s is required", v.Key))
			}
		}
	default:
		return "", xerror.ValidationFailed(nil, fmt.Sprintf("Unsupported trigger type: %s", triggerData.TriggerType))
	}

	// Create workflow execution
	wfe := workflowexecution.NewWorkflowExecution(wf.ID, wf.Data, params.RuntimeVariables)

	// Add trigger node to steps
	var steps []stepexecution.StepExecution
	triggerStep := stepexecution.NewStepExecution(wfe.ID, *triggerNode, params.RuntimeVariables)
	steps = append(steps, triggerStep)

	// Add other nodes to steps
	for _, n := range wf.Data.Nodes {
		if n.Type == node.TypeTrigger {
			continue
		}
		step := stepexecution.NewStepExecution(wfe.ID, n, params.RuntimeVariables)
		steps = append(steps, step)
	}

	if err = s.sqlDBProvider.WithTx(ctx, func(db sqldb.SQLDB) error {
		if err := s.workflowExecutionRepo.CreateWorkflowExecution(ctx, db, wfe); err != nil {
			return fmt.Errorf("repo create workflow execution: %w", err)
		}

		if err := s.stepExecutionRepo.BatchCreateStepExecutions(ctx, db, steps); err != nil {
			return fmt.Errorf("repo batch create step executions: %w", err)
		}

		return nil
	}); err != nil {
		return "", fmt.Errorf("with tx: %w", err)
	}

	// Publish event
	ev := pubsub.WorkflowExecutionCreated{
		WorkflowExecutionID: wfe.ID,
	}
	payload, err := json.Marshal(ev)
	if err != nil {
		return "", fmt.Errorf("marshal event: %w", err)
	}

	msg := message.NewMessage(uuid.NewString(), payload)
	if err := s.publisher.Publish(pubsub.WorkflowExecutionCreatedTopic, msg); err != nil {
		return "", fmt.Errorf("publisher publish event: %w", err)
	}

	return wfe.ID, nil
}
