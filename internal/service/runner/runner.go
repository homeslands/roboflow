package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/tuanvumaihuynh/roboflow/internal/model"
	"github.com/tuanvumaihuynh/roboflow/internal/service/raybot"
	raybotcommand "github.com/tuanvumaihuynh/roboflow/internal/service/raybot_command"
	workflowexecution "github.com/tuanvumaihuynh/roboflow/internal/service/workflow_execution"
)

type Runner struct {
	stepRepo                  model.StepRepository
	raybotSvc                 raybot.Service
	raybotCommandSvc          raybotcommand.Service
	workflowExecutionSvc      workflowexecution.Service
	raybotCommandStatusPoller RaybotCommandStatusPoller

	log *slog.Logger
}

func NewRunner(
	stepRepo model.StepRepository,
	raybotCommandRepo model.RaybotCommandRepository,
	raybotSvc raybot.Service,
	raybotCommandSvc raybotcommand.Service,
	workflowExecutionSvc workflowexecution.Service,
	log *slog.Logger,
) *Runner {
	return &Runner{
		stepRepo:             stepRepo,
		raybotSvc:            raybotSvc,
		raybotCommandSvc:     raybotCommandSvc,
		workflowExecutionSvc: workflowExecutionSvc,
		raybotCommandStatusPoller: RaybotCommandStatusPoller{
			raybotCommandRepo: raybotCommandRepo,
		},
		log: log,
	}
}

func (r *Runner) HandleRunWorkflowExecution(msg *message.Message) (err error) {
	var wfe model.WorkflowExecution
	if err := json.Unmarshal(msg.Payload, &wfe); err != nil {
		r.log.Error("Failed to unmarshal message", slog.Any("error", err))
		return nil
	}

	defer func() {
		if err != nil {
			r.log.Error("Failed to handle run workflow execution", slog.Any("error", err))
			err = r.workflowExecutionSvc.SetFailed(msg.Context(), workflowexecution.SetWorkflowExecutionFailedCommand{
				ID: wfe.ID,
			})
			if err != nil {
				r.log.Error("Failed to set workflow execution status to FAILED", slog.Any("error", err))
			}
		}
	}()

	r.log.Info("Start running workflow execution",
		slog.String("workflow_execution_id", wfe.ID.String()))

	//  Set workflow execution status to RUNNING
	err = r.workflowExecutionSvc.SetRunning(msg.Context(), workflowexecution.SetWorkflowExecutionRunningCommand{
		ID: wfe.ID,
	})
	if err != nil {
		r.log.Error("Failed to set workflow execution status to RUNNING", slog.Any("error", err))
		return
	}

	// Run steps
	phases, err := generatePhases(wfe)
	if err != nil {
		r.log.Error("Failed to generate phases", slog.Any("error", err))
		return
	}

	for _, phase := range phases {
		for _, step := range phase.Steps {
			r.log.Info("Executing step", slog.String("step_id", step.ID.String()))
			err = r.executeStep(msg.Context(), step, *r)
			if err != nil {
				r.log.Error("Failed to execute step",
					slog.String("step_id", step.ID.String()),
					slog.Any("error", err))
				return
			}
		}
	}

	// Set workflow execution status to COMPLETED
	err = r.workflowExecutionSvc.SetCompleted(msg.Context(), workflowexecution.SetWorkflowExecutionCompletedCommand{
		ID: wfe.ID,
	})
	if err != nil {
		r.log.Error("Failed to set workflow execution status to COMPLETED", slog.Any("error", err))
	}

	return
}

func (r *Runner) executeStep(ctx context.Context, step model.Step, runner Runner) (err error) {
	handler, err := getTaskHandler(step, runner)
	if err != nil {
		return fmt.Errorf("failed to get task handler: %w", err)
	}

	defer func() {
		// Post step execution
		now := time.Now()
		step.CompletedAt = &now
		if err != nil {
			step.Status = model.WorkflowExecutionStepStatusFailed
		} else {
			step.Status = model.WorkflowExecutionStepStatusCompleted
		}
		updateErr := r.stepRepo.Update(ctx, step)
		if updateErr != nil {
			err = fmt.Errorf("failed to update step: %w", updateErr)
		}
	}()

	// Pre step execution
	now := time.Now()
	step.StartedAt = &now
	step.Status = model.WorkflowExecutionStepStatusRunning
	err = r.stepRepo.Update(ctx, step)
	if err != nil {
		return fmt.Errorf("failed to update step: %w", err)
	}

	// Execute step
	err = handler.Handle(ctx, step)
	if err != nil {
		return fmt.Errorf("failed to execute step: %w", err)
	}

	return err
}

func getTaskHandler(step model.Step, runner Runner) (TaskHandler, error) {
	switch step.Node.Definition.Type {
	case model.TaskRaybotValidateState:
		return RaybotValidateStateTaskHandler{
			raybotSvc: runner.raybotSvc,
		}, nil
	case model.TaskRaybotMoveToLocation:
		return RaybotMoveToLocationTaskHandler{
			poller:           runner.raybotCommandStatusPoller,
			raybotCommandSvc: runner.raybotCommandSvc,
		}, nil
	case model.TaskRaybotMoveForward:
		return RaybotMoveForwardTaskHandler{
			poller:           runner.raybotCommandStatusPoller,
			raybotCommandSvc: runner.raybotCommandSvc,
		}, nil
	case model.TaskRaybotMoveBackward:
		return RaybotMoveBackwardTaskHandler{
			poller:           runner.raybotCommandStatusPoller,
			raybotCommandSvc: runner.raybotCommandSvc,
		}, nil
	case model.TaskRaybotOpenBox:
		return RaybotOpenBoxTaskHandler{
			poller:           runner.raybotCommandStatusPoller,
			raybotCommandSvc: runner.raybotCommandSvc,
		}, nil
	case model.TaskRaybotCloseBox:
		return RaybotCloseBoxTaskHandler{
			poller:           runner.raybotCommandStatusPoller,
			raybotCommandSvc: runner.raybotCommandSvc,
		}, nil
	case model.TaskRaybotLiftBox:
		return RaybotLiftBoxTaskHandler{
			poller:           runner.raybotCommandStatusPoller,
			raybotCommandSvc: runner.raybotCommandSvc,
		}, nil
	case model.TaskRaybotDropBox:
		return RaybotDropBoxTaskHandler{
			poller:           runner.raybotCommandStatusPoller,
			raybotCommandSvc: runner.raybotCommandSvc,
		}, nil
	case model.TaskRaybotCheckQRCode:
		return RaybotCheckQRCodeTaskHandler{
			poller:           runner.raybotCommandStatusPoller,
			raybotCommandSvc: runner.raybotCommandSvc,
		}, nil
	case model.TaskRaybotWaitGetItem:
		return RaybotWaitGetItemTaskHandler{
			poller:           runner.raybotCommandStatusPoller,
			raybotCommandSvc: runner.raybotCommandSvc,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported step type: %s", step.Node.Definition.Type)
	}
}
