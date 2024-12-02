package client

import (
	"encoding/json"

	cmdModel "github.com/tuanvumaihuynh/roboflow/internal/command/model"
	raybotModel "github.com/tuanvumaihuynh/roboflow/internal/raybot/model"
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
	ID   string                 `json:"id"`
	Type cmdModel.CommandType   `json:"type"`
	Data map[string]interface{} `json:"data"`
}

// Publish event data

type PublishUpdateStatus struct {
	Status raybotModel.RaybotStatus `json:"status"`
}

type PublishUpdateRamCpu struct {
	RamUsage float64 `json:"ram_usage"`
	CpuUsage float64 `json:"cpu_usage"`
}

// Response event data

type ResponseStatus struct {
	Status raybotModel.RaybotStatus `json:"status"`
}
