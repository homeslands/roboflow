package raybotclient

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
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
	Topic InboundTopic    `json:"topic"`
	Data  json.RawMessage `json:"data"`
}

func (m InboundPublishMsg) Operation() Operation { return OperationPublish }

func (m *InboundPublishMsg) UnmarshalJSON(data []byte) error {
	type Alias InboundPublishMsg
	alias := (*Alias)(m)
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	// Validate the topic field
	if err := alias.Topic.Validate(); err != nil {
		return err
	}

	return nil
}

type InboundTopic string

func (t InboundTopic) Validate() error {
	if t == "" {
		return fmt.Errorf("topic is empty")
	}
	switch t {
	case TopicStatus,
		TopicRamCpu,
		TopicLog,
		TopicBatterySensor,
		TopicWeightSensor,
		TopicForwardDistanceSensor,
		TopicBackwardDistanceSensor,
		TopicMovementMotor,
		TopicLiftingMotor:
	default:
		return fmt.Errorf("invalid topic: %s", t)
	}
	return nil
}

const (
	TopicStatus                 InboundTopic = "status"
	TopicRamCpu                 InboundTopic = "ram_cpu"
	TopicLog                    InboundTopic = "log"
	TopicBatterySensor          InboundTopic = "battery_sensor"
	TopicWeightSensor           InboundTopic = "weight_sensor"
	TopicForwardDistanceSensor  InboundTopic = "forward_distance_sensor"
	TopicBackwardDistanceSensor InboundTopic = "backward_distance_sensor"
	TopicMovementMotor          InboundTopic = "movement_motor"
	TopicLiftingMotor           InboundTopic = "lifting_motor"
)

type CommandStatus string

func (s CommandStatus) Validate() error {
	switch s {
	case CommandStatusInProgress,
		CommandStatusSuccess,
		CommandStatusError:
	default:
		return fmt.Errorf("invalid command status: %s", s)
	}
	return nil
}

const (
	CommandStatusInProgress CommandStatus = "IN_PROGRESS"
	CommandStatusSuccess    CommandStatus = "SUCCESS"
	CommandStatusError      CommandStatus = "ERROR"
)

// InboundResponseMsg is the response message received from the client
//
// Now InboundResponseMsg only use for command response so
// in the struct we have status field to represent the status of the command
type InboundResponseMsg struct {
	ID     uuid.UUID       `json:"id"`
	Data   json.RawMessage `json:"data"`
	Status CommandStatus   `json:"status"`
}

func (m InboundResponseMsg) Operation() Operation { return OperationResponse }

func (m *InboundResponseMsg) UnmarshalJSON(data []byte) error {
	type Alias InboundResponseMsg
	alias := (*Alias)(m)
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	// Validate the status field
	if err := alias.Status.Validate(); err != nil {
		return err
	}

	return nil
}
