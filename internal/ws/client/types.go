package client

import (
	"encoding/json"

	"github.com/tuanvumaihuynh/roboflow/internal/model"
)

type Operation string

const (
	OperationPublish  Operation = "publish"
	OperationResponse Operation = "response"
)

type InboundMsg struct {
	Operation Operation       `json:"op"`
	Id        *string         `json:"id,omitempty"`
	Topic     *string         `json:"topic,omitempty"`
	Data      json.RawMessage `json:"data"`
}

// OutboundCommandMsg is the message sent to the client when a command is received
type OutboundCommandMsg struct {
	ID   string                  `json:"id"`
	Type model.RaybotCommandType `json:"type"`
	Data map[string]interface{}  `json:"data"`
}

// Publish event data

type PublishUpdateStatus struct {
	Status model.RaybotStatus `json:"status"`
}

type PublishUpdateRamCpu struct {
	RamUsage float64 `json:"ram_usage"`
	CpuUsage float64 `json:"cpu_usage"`
}

// Response event data

type ResponseStatus struct {
	Status model.RaybotStatus `json:"status"`
}
