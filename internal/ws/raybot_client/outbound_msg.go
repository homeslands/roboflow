package raybotclient

import (
	"encoding/json"

	"github.com/tuanvumaihuynh/roboflow/internal/model"
)

// OutboundCommandMsg represents a command message that is sent to the raybot.
type OutboundCommandMsg struct {
	ID   string                  `json:"id"`
	Type model.RaybotCommandType `json:"type"`
	Data json.RawMessage         `json:"data,omitempty"`
}
