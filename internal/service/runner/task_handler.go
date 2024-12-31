package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/tuanvumaihuynh/roboflow/internal/model"
	"github.com/tuanvumaihuynh/roboflow/internal/service/raybot"
	raybotcommand "github.com/tuanvumaihuynh/roboflow/internal/service/raybot_command"
)

type TaskHandler interface {
	Handle(ctx context.Context, step model.Step) error
}

type RaybotValidateStateTaskHandler struct {
	raybotSvc raybot.Service
}

func (e RaybotValidateStateTaskHandler) Handle(ctx context.Context, step model.Step) error {
	raybotIDField := step.Node.Definition.Fields[model.NodeDefinitionFieldRaybotID]
	raybotID := raybotIDField.Value
	if raybotID == nil {
		return fmt.Errorf("raybot id is required")
	}

	r, err := e.raybotSvc.GetByID(ctx, raybot.GetRaybotByIDQuery{
		ID: uuid.MustParse(*raybotID),
	})
	if err != nil {
		return fmt.Errorf("failed to get raybot by id: %w", err)
	}

	if r.Status != model.RaybotStatusIdle {
		return fmt.Errorf("raybot is not idle, current status: %s", r.Status)
	}

	return nil
}

type RaybotMoveToLocationTaskHandler struct {
	poller           RaybotCommandStatusPoller
	raybotCommandSvc raybotcommand.Service
}

func (e RaybotMoveToLocationTaskHandler) Handle(ctx context.Context, step model.Step) error {
	// Prepare input
	raybotIDField := step.Node.Definition.Fields[model.NodeDefinitionFieldRaybotID]
	raybotID := raybotIDField.Value
	if raybotID == nil {
		return fmt.Errorf("raybot id is required")
	}

	directionField := step.Node.Definition.Fields[model.NodeDefinitionFieldDirection]
	direction := directionField.Value
	if direction == nil {
		return fmt.Errorf("direction is required")
	}

	locationField := step.Node.Definition.Fields[model.NodeDefinitionFieldLocation]
	location := locationField.Value
	if location == nil {
		return fmt.Errorf("location is required")
	}

	// Create raybot command
	input := model.MoveToLocationInput{
		Location:  *location,
		Direction: *direction,
	}
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("failed to marshal input: %w", err)
	}
	cmd, err := e.raybotCommandSvc.Create(ctx, raybotcommand.CreateRaybotCommandCommand{
		RaybotID: uuid.MustParse(*raybotID),
		Type:     model.RaybotCommandTypeMoveToLocation,
		Input:    inputBytes,
	})
	if err != nil {
		return fmt.Errorf("failed to create raybot command: %w", err)
	}

	timeOutSec := step.Node.Definition.TimeoutSec
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeOutSec)*time.Second)
	defer cancel()
	return e.poller.pollRaybotCommandStatus(ctx, cmd.ID)
}

type RaybotMoveForwardTaskHandler struct {
	poller           RaybotCommandStatusPoller
	raybotCommandSvc raybotcommand.Service
}

func (e RaybotMoveForwardTaskHandler) Handle(ctx context.Context, step model.Step) error {
	// Prepare input
	raybotIDField := step.Node.Definition.Fields[model.NodeDefinitionFieldRaybotID]
	raybotID := raybotIDField.Value
	if raybotID == nil {
		return fmt.Errorf("raybot id is required")
	}

	// Create raybot command
	cmd, err := e.raybotCommandSvc.Create(ctx, raybotcommand.CreateRaybotCommandCommand{
		RaybotID: uuid.MustParse(*raybotID),
		Type:     model.RaybotCommandTypeMoveForward,
	})
	if err != nil {
		return fmt.Errorf("failed to create raybot command: %w", err)
	}

	timeOutSec := step.Node.Definition.TimeoutSec
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeOutSec)*time.Second)
	defer cancel()
	return e.poller.pollRaybotCommandStatus(ctx, cmd.ID)
}

type RaybotMoveBackwardTaskHandler struct {
	poller           RaybotCommandStatusPoller
	raybotCommandSvc raybotcommand.Service
}

func (e RaybotMoveBackwardTaskHandler) Handle(ctx context.Context, step model.Step) error {
	// Prepare input
	raybotIDField := step.Node.Definition.Fields[model.NodeDefinitionFieldRaybotID]
	raybotID := raybotIDField.Value
	if raybotID == nil {
		return fmt.Errorf("raybot id is required")
	}

	// Create raybot command
	cmd, err := e.raybotCommandSvc.Create(ctx, raybotcommand.CreateRaybotCommandCommand{
		RaybotID: uuid.MustParse(*raybotID),
		Type:     model.RaybotCommandTypeMoveBackward,
	})
	if err != nil {
		return fmt.Errorf("failed to create raybot command: %w", err)
	}

	timeOutSec := step.Node.Definition.TimeoutSec
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeOutSec)*time.Second)
	defer cancel()
	return e.poller.pollRaybotCommandStatus(ctx, cmd.ID)
}

type RaybotOpenBoxTaskHandler struct {
	poller           RaybotCommandStatusPoller
	raybotCommandSvc raybotcommand.Service
}

func (e RaybotOpenBoxTaskHandler) Handle(ctx context.Context, step model.Step) error {
	// Prepare input
	raybotIDField := step.Node.Definition.Fields[model.NodeDefinitionFieldRaybotID]
	raybotID := raybotIDField.Value
	if raybotID == nil {
		return fmt.Errorf("raybot id is required")
	}

	// Create raybot command
	cmd, err := e.raybotCommandSvc.Create(ctx, raybotcommand.CreateRaybotCommandCommand{
		RaybotID: uuid.MustParse(*raybotID),
		Type:     model.RaybotCommandTypeOpenBox,
	})
	if err != nil {
		return fmt.Errorf("failed to create raybot command: %w", err)
	}

	timeOutSec := step.Node.Definition.TimeoutSec
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeOutSec)*time.Second)
	defer cancel()
	return e.poller.pollRaybotCommandStatus(ctx, cmd.ID)
}

type RaybotCloseBoxTaskHandler struct {
	poller           RaybotCommandStatusPoller
	raybotCommandSvc raybotcommand.Service
}

func (e RaybotCloseBoxTaskHandler) Handle(ctx context.Context, step model.Step) error {
	// Prepare input
	raybotIDField := step.Node.Definition.Fields[model.NodeDefinitionFieldRaybotID]
	raybotID := raybotIDField.Value
	if raybotID == nil {
		return fmt.Errorf("raybot id is required")
	}

	// Create raybot command
	cmd, err := e.raybotCommandSvc.Create(ctx, raybotcommand.CreateRaybotCommandCommand{
		RaybotID: uuid.MustParse(*raybotID),
		Type:     model.RaybotCommandTypeCloseBox,
	})
	if err != nil {
		return fmt.Errorf("failed to create raybot command: %w", err)
	}

	timeOutSec := step.Node.Definition.TimeoutSec
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeOutSec)*time.Second)
	defer cancel()
	return e.poller.pollRaybotCommandStatus(ctx, cmd.ID)
}

type RaybotLiftBoxTaskHandler struct {
	poller           RaybotCommandStatusPoller
	raybotCommandSvc raybotcommand.Service
}

func (e RaybotLiftBoxTaskHandler) Handle(ctx context.Context, step model.Step) error {
	// Prepare input
	raybotIDField := step.Node.Definition.Fields[model.NodeDefinitionFieldRaybotID]
	raybotID := raybotIDField.Value
	if raybotID == nil {
		return fmt.Errorf("raybot id is required")
	}

	// Create raybot command
	cmd, err := e.raybotCommandSvc.Create(ctx, raybotcommand.CreateRaybotCommandCommand{
		RaybotID: uuid.MustParse(*raybotID),
		Type:     model.RaybotCommandTypeLiftBox,
	})
	if err != nil {
		return fmt.Errorf("failed to create raybot command: %w", err)
	}

	timeOutSec := step.Node.Definition.TimeoutSec
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeOutSec)*time.Second)
	defer cancel()
	return e.poller.pollRaybotCommandStatus(ctx, cmd.ID)
}

type RaybotDropBoxTaskHandler struct {
	poller           RaybotCommandStatusPoller
	raybotCommandSvc raybotcommand.Service
}

func (e RaybotDropBoxTaskHandler) Handle(ctx context.Context, step model.Step) error {
	// Prepare input
	raybotIDField := step.Node.Definition.Fields[model.NodeDefinitionFieldRaybotID]
	raybotID := raybotIDField.Value
	if raybotID == nil {
		return fmt.Errorf("raybot id is required")
	}

	// Create raybot command
	cmd, err := e.raybotCommandSvc.Create(ctx, raybotcommand.CreateRaybotCommandCommand{
		RaybotID: uuid.MustParse(*raybotID),
		Type:     model.RaybotCommandTypeDropBox,
	})
	if err != nil {
		return fmt.Errorf("failed to create raybot command: %w", err)
	}

	timeOutSec := step.Node.Definition.TimeoutSec
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeOutSec)*time.Second)
	defer cancel()
	return e.poller.pollRaybotCommandStatus(ctx, cmd.ID)
}

type RaybotCheckQRCodeTaskHandler struct {
	poller           RaybotCommandStatusPoller
	raybotCommandSvc raybotcommand.Service
}

func (e RaybotCheckQRCodeTaskHandler) Handle(ctx context.Context, step model.Step) error {
	// Prepare input
	raybotIDField := step.Node.Definition.Fields[model.NodeDefinitionFieldRaybotID]
	raybotID := raybotIDField.Value
	if raybotID == nil {
		return fmt.Errorf("raybot id is required")
	}

	// Create raybot command
	cmd, err := e.raybotCommandSvc.Create(ctx, raybotcommand.CreateRaybotCommandCommand{
		RaybotID: uuid.MustParse(*raybotID),
		Type:     model.RaybotCommandTypeCheckQrCode,
	})
	if err != nil {
		return fmt.Errorf("failed to create raybot command: %w", err)
	}

	timeOutSec := step.Node.Definition.TimeoutSec
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeOutSec)*time.Second)
	defer cancel()
	return e.poller.pollRaybotCommandStatus(ctx, cmd.ID)
}

type RaybotWaitGetItemTaskHandler struct {
	poller           RaybotCommandStatusPoller
	raybotCommandSvc raybotcommand.Service
}

func (e RaybotWaitGetItemTaskHandler) Handle(ctx context.Context, step model.Step) error {
	// Prepare input
	raybotIDField := step.Node.Definition.Fields[model.NodeDefinitionFieldRaybotID]
	raybotID := raybotIDField.Value
	if raybotID == nil {
		return fmt.Errorf("raybot id is required")
	}

	// Create raybot command
	cmd, err := e.raybotCommandSvc.Create(ctx, raybotcommand.CreateRaybotCommandCommand{
		RaybotID: uuid.MustParse(*raybotID),
		Type:     model.RaybotCommandTypeWaitGetItem,
	})
	if err != nil {
		return fmt.Errorf("failed to create raybot command: %w", err)
	}

	timeOutSec := step.Node.Definition.TimeoutSec
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeOutSec)*time.Second)
	defer cancel()
	return e.poller.pollRaybotCommandStatus(ctx, cmd.ID)
}
