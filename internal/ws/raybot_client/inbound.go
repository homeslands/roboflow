package raybotclient

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/tuanvumaihuynh/roboflow/internal/model"
)

type Operation string

const (
	OperationPublish  Operation = "publish"
	OperationResponse Operation = "response"
)

// InboundMsg is the message received from the client
type InboundMsg interface {
	Operation() Operation
}

type InboundPublishMsg struct {
	Topic model.RaybotCommandType `json:"topic"`
	Data  json.RawMessage         `json:"data"`
}

func (m InboundPublishMsg) Operation() Operation { return OperationPublish }

type InboundResponseMsg struct {
	ID   uuid.UUID       `json:"id"`
	Data json.RawMessage `json:"data"`
}

func (m InboundResponseMsg) Operation() Operation { return OperationResponse }

// UnmarshalInboundMsg unmarshals the inbound message from raybot robot client.
// The operation field is used to determine the type of the message.
func UnmarshalInboundMsg(data []byte) (InboundMsg, error) {
	var msg InboundMsg
	var temp struct {
		Op Operation `json:"op"`
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return nil, err
	}

	switch temp.Op {
	case OperationPublish:
		msg = &InboundPublishMsg{}
	case OperationResponse:
		msg = &InboundResponseMsg{}
	default:
		return nil, fmt.Errorf("invalid operation: %s", temp.Op)
	}
	if err := json.Unmarshal(data, msg); err != nil {
		return nil, err
	}
	return msg, nil
}
