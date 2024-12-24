package hub

import (
	"encoding/json"
	"log/slog"

	"github.com/ThreeDotsLabs/watermill/message"

	raybotCommandEvent "github.com/tuanvumaihuynh/roboflow/internal/service/raybot_command/event"
)

func (h *Hub) HandleCommandCreated(msg *message.Message) error {
	var cmdCreated raybotCommandEvent.RaybotCommandCreated
	err := json.Unmarshal(msg.Payload, &cmdCreated)
	if err != nil {
		h.logger.Error("Error unmarshalling command created event", slog.Any("error", err))
		return nil
	}

	h.logger.Info("Command created event received", slog.String("id", cmdCreated.ID.String()))

	// Find the client
	client, exists := h.clients[cmdCreated.RaybotID.String()]
	if !exists {
		h.logger.Error("Raybot client not found", slog.String("id", cmdCreated.RaybotID.String()))
		return nil
	}

	// Send the command to the client
	h.logger.Info("Sending command to client", slog.String("id", cmdCreated.ID.String()))
	client.SendCommand(&cmdCreated)
	return nil
}
