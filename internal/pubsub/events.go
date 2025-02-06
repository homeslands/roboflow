package pubsub

import raybotcommand "github.com/tuanvumaihuynh/roboflow/internal/model/raybot_command"

const (
	RaybotCommandCreatedTopic     = "raybot_command:created"
	WorkflowExecutionCreatedTopic = "workflow_execution:created"
)

type RaybotCommandCreated struct {
	RaybotID  string               `json:"raybot_id"`
	CommandID string               `json:"command_id"`
	Type      raybotcommand.Type   `json:"type"`
	Inputs    raybotcommand.Inputs `json:"inputs"`
}

type WorkflowExecutionCreated struct {
	WorkflowExecutionID string `json:"workflow_execution_id"`
}
