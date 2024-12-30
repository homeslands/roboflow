package raybotclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/gorilla/websocket"

	"github.com/tuanvumaihuynh/roboflow/internal/model"
	raybotcommand "github.com/tuanvumaihuynh/roboflow/internal/service/raybot_command"
	"github.com/tuanvumaihuynh/roboflow/pkg/xerrors"
)

func handleInboundResponseMsg(c *RaybotClient, msg InboundResponseMsg) {
	c.log.Debug("Received response msg",
		slog.String("cmd_id", msg.ID.String()),
		slog.Any("data", msg.Data))

	// Get current command
	cmd := c.GetCurrentCmd()
	if cmd == nil {
		c.log.Error("Raybot is not processing any command")
		closeConn(c, websocket.CloseInvalidFramePayloadData, "raybot is not processing any command")
		return
	}

	// Check duplicate response
	if msg.Status == CommandStatusInProgress &&
		cmd.Status == model.RaybotCommandStatusInProgress {
		c.log.Warn("Duplicate response")
		return
	}

	switch msg.Status {
	case CommandStatusInProgress:
		handleInProgress(c, msg, *cmd)
	case CommandStatusSuccess:
		handleSuccess(c, msg, *cmd)
	case CommandStatusError:
		handleError(c, msg, *cmd)
	}
}

// handleInProgress handles in progress response from the raybot.
func handleInProgress(c *RaybotClient, msg InboundResponseMsg, currentCmd model.RaybotCommand) {
	c.log.Info("Handle command IN_PROGRESS",
		slog.String("cmd_id", msg.ID.String()))

	// Update command
	err := c.raybotCommandSvc.SetStatusInProgess(context.TODO(), raybotcommand.SetStatusInProgessCommand{
		ID: currentCmd.ID,
	})
	if err != nil {
		c.log.Error("Error set command status IN_PROGRESS", slog.Any("error", err))
		closeConn(c, websocket.CloseInternalServerErr, "internal server error")
		return
	}
}

// handleSuccess handles success response from the raybot.
func handleSuccess(c *RaybotClient, msg InboundResponseMsg, currentCmd model.RaybotCommand) {
	c.log.Info("Handle command SUCCESS",
		slog.String("cmd_id", msg.ID.String()))

	// Set current raybot command to nil
	c.SetCurrentCmd(nil)

	output, err := getCommandOutput(currentCmd.Type, msg)
	if err != nil {
		c.log.Error("Error getting command output", slog.Any("error", err))
		closeConn(c, websocket.CloseInvalidFramePayloadData, "invalid command output for success response")
		go c.markCommandAsFailed(context.TODO(), currentCmd, model.FailedOutput{
			Reason: "invalid command output for success response",
		})
		return
	}

	err = c.raybotCommandSvc.SetStatusSuccess(context.TODO(), raybotcommand.SetStatusSuccessCommand{
		ID:     currentCmd.ID,
		Output: output,
	})
	if err != nil {
		c.log.Error("Error set command status SUCCESS", slog.Any("error", err))
		closeConn(c, websocket.CloseInternalServerErr, "internal server error")
		return
	}
}

// handleError handles error response from the raybot.
func handleError(c *RaybotClient, msg InboundResponseMsg, currentCmd model.RaybotCommand) {
	c.log.Info("Handle command FAILED",
		slog.String("cmd_id", msg.ID.String()))

	var data InboundResponseErrorData
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		c.log.Error("Error parsing error data",
			slog.Any("error", err),
			slog.String("raw_data", string(msg.Data)))

		closeConn(c, websocket.CloseInvalidFramePayloadData, "invalid error data")
		go c.markCommandAsFailed(context.TODO(), currentCmd, model.FailedOutput{
			Reason: "invalid error data",
		})
		return
	}

	// Set current raybot command to nil
	c.SetCurrentCmd(nil)

	c.markCommandAsFailed(context.TODO(), currentCmd, model.FailedOutput{
		Reason: data.Reason,
	})
}

func (c *RaybotClient) markCommandAsFailed(ctx context.Context, cmd model.RaybotCommand, output model.FailedOutput) {
	c.log.Error("Marking command as failed",
		slog.String("cmd_id", cmd.ID.String()),
		slog.String("reason", output.Reason))

	outputBytes, err := json.Marshal(output)
	if err != nil {
		c.log.Error("Failed to marshal error output", slog.Any("error", err))
		return
	}

	err = c.raybotCommandSvc.SetStatusFailed(ctx, raybotcommand.SetStatusFailedCommand{
		ID:     cmd.ID,
		Output: outputBytes,
	})
	if err != nil {
		if xerrors.IsNotFound(err) {
			c.log.Warn("Command not found", slog.String("cmd_id", cmd.ID.String()))
			return
		}
		c.log.Error("Failed to update command status to FAILED", slog.Any("error", err))
	}
}

// getCommandOutput processes the output based on the command type.
func getCommandOutput(cmdType model.RaybotCommandType, msg InboundResponseMsg) ([]byte, error) {
	switch cmdType {
	case model.RaybotCommandTypeStop,
		model.RaybotCommandTypeMoveForward,
		model.RaybotCommandTypeMoveBackward,
		model.RaybotCommandTypeMoveToLocation,
		model.RaybotCommandTypeOpenBox,
		model.RaybotCommandTypeCloseBox,
		model.RaybotCommandTypeLiftBox,
		model.RaybotCommandTypeDropBox,
		model.RaybotCommandTypeCheckQrCode,
		model.RaybotCommandTypeWaitGetItem:
		return nil, nil
	case model.RaybotCommandTypeScanLocation:
		var data InboundResponseScanLocationData
		if err := json.Unmarshal(msg.Data, &data); err != nil {
			return nil, fmt.Errorf("unmarshal scan location data: %w", err)
		}
		output := model.ScanLocationOutput{
			Locations: data.Locations,
		}
		outputBytes, err := json.Marshal(output)
		if err != nil {
			return nil, fmt.Errorf("marshal scan location output: %w", err)
		}

		return outputBytes, nil

	default:
		return nil, errors.New("unsupported command type")
	}
}
