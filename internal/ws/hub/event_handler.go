package hub

import (
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill/message"
	"go.uber.org/zap"

	cmdEvent "github.com/tuanvumaihuynh/roboflow/internal/command/event"
)

func (h *Hub) HandleCommandCreated(msg *message.Message) error {
	var cmdCreated cmdEvent.CommandCreated
	err := json.Unmarshal(msg.Payload, &cmdCreated)
	if err != nil {
		h.logger.Error("Error unmarshalling command created event", zap.Error(err))
		return nil
	}

	h.logger.Info("Command created event received", zap.String("id", cmdCreated.ID.String()))

	// Find the client
	client, exists := h.clients[cmdCreated.RaybotID.String()]
	if !exists {
		h.logger.Error("Raybot client not found", zap.String("id", cmdCreated.RaybotID.String()))
		return nil
	}

	// Send the command to the client
	h.logger.Info("Sending command to client", zap.String("id", cmdCreated.ID.String()))
	client.SendCommand(&cmdCreated)
	return nil
}
