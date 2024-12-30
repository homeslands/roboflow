package raybotclient

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/tuanvumaihuynh/roboflow/internal/model"
	raybotcommand "github.com/tuanvumaihuynh/roboflow/internal/service/raybot_command"
	"github.com/tuanvumaihuynh/roboflow/internal/service/raybot_command/event"
)

func (ws *WebSocket) HandleRaybotCommandCreated(msg []byte) error {
	var cmd event.RaybotCommandCreated
	if err := json.Unmarshal(msg, &cmd); err != nil {
		ws.log.Error("Failed to unmarshal raybot command created event", slog.Any("error", err))
		return nil
	}
	ws.log.Debug("Raybot command created event received",
		slog.String("command_id", cmd.ID.String()),
		slog.String("raybot_id", cmd.RaybotID.String()))

	// Get raybot client
	c := ws.getClient(cmd.RaybotID)
	if c == nil {
		ws.log.Info("Raybot is not connected",
			slog.String("command_id", cmd.ID.String()),
			slog.String("raybot_id", cmd.RaybotID.String()))

		// Handle unprocessed raybot command
		// Publish event
		go ws.handleUnprocessedRaybotCommandCreated(cmd, model.FailedOutput{
			Reason: "Raybot is not connected",
		})
		return nil
	}

	// Set current command
	c.SetCurrentCmd(&cmd)

	// Create outbound command message
	outboundCmdMsg := OutboundCommandMsg{
		ID:   cmd.ID.String(),
		Type: cmd.Type,
		Data: cmd.Input,
	}
	msgBytes, err := json.Marshal(outboundCmdMsg)
	if err != nil {
		ws.log.Error("Failed to marshal outbound command message", slog.Any("error", err))
		return nil
	}

	// Send command to raybot
	select {
	case c.outboundChan <- msgBytes:
		c.log.Debug("Command sent to raybot",
			slog.String("command_id", cmd.ID.String()),
			slog.String("command_type", string(cmd.Type)))
	default:
		c.log.Error("Outbound channel is full, drop message",
			slog.String("command_id", cmd.ID.String()),
			slog.Int("channel_size", len(c.outboundChan)))
	}

	return nil
}

func (ws *WebSocket) handleUnprocessedRaybotCommandCreated(cmd model.RaybotCommand, output model.FailedOutput) {
	ws.log.Info("Unprocessed raybot command created event",
		slog.String("command_id", cmd.ID.String()),
		slog.String("raybot_id", cmd.RaybotID.String()))

	outputBytes, err := json.Marshal(output)
	if err != nil {
		ws.log.Error("Failed to marshal output", slog.Any("error", err))
		return
	}

	err = ws.raybotCommandSvc.SetStatusFailed(context.TODO(), raybotcommand.SetStatusFailedCommand{
		ID:     cmd.ID,
		Output: outputBytes,
	})
	if err != nil {
		ws.log.Error("Failed to set raybot command status failed", slog.Any("error", err))
	}
}

func getClientIPAddress(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}

	return ip
}
