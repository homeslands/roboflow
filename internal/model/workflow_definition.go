package model

import "github.com/tuanvumaihuynh/roboflow/pkg/xerrors"

type TaskType string

func (t TaskType) Validate() error {
	switch t {
	case TaskRaybotValidateState:
	case TaskRaybotMoveForward:
	case TaskRaybotMoveBackward:
	case TaskRaybotMoveToLocation:
	case TaskRaybotOpenBox:
	case TaskRaybotCloseBox:
	case TaskRaybotLiftBox:
	case TaskRaybotDropBox:
	case TaskRaybotCheckQRCode:
	case TaskRaybotWaitGetItem:
	default:
		return xerrors.ThrowInvalidArgument(nil, "invalid task type")
	}
	return nil
}

const (
	TaskRaybotValidateState TaskType = "VALIDATE_STATE"

	TaskRaybotMoveForward    TaskType = "MOVE_FORWARD"
	TaskRaybotMoveBackward   TaskType = "MOVE_BACKWARD"
	TaskRaybotMoveToLocation TaskType = "MOVE_TO_LOCATION"

	TaskRaybotOpenBox  TaskType = "OPEN_BOX"
	TaskRaybotCloseBox TaskType = "CLOSE_BOX"
	TaskRaybotLiftBox  TaskType = "LIFT_BOX"
	TaskRaybotDropBox  TaskType = "DROP_BOX"

	TaskRaybotCheckQRCode TaskType = "CHECK_QR"

	TaskRaybotWaitGetItem TaskType = "WAIT_GET_ITEM"
)

type NodeType string

const (
	NodeTypeRaybot NodeType = "RaybotNode"
)

func (n NodeType) Validate() error {
	switch n {
	case NodeTypeRaybot:
	default:
		return xerrors.ThrowInvalidArgument(nil, "invalid node type")
	}
	return nil
}

type WorkflowEdge struct {
	ID           string  `json:"id" validate:"required"`
	Type         string  `json:"type" validate:"required"`
	Source       string  `json:"source" validate:"required"`
	Target       string  `json:"target" validate:"required"`
	SourceHandle string  `json:"source_handle" validate:"required"`
	TargetHandle string  `json:"target_handle" validate:"required"`
	Label        string  `json:"label" validate:"required"`
	Animated     bool    `json:"animated" validate:"required"`
	SourceX      float32 `json:"source_x" validate:"required"`
	SourceY      float32 `json:"source_y" validate:"required"`
	TargetX      float32 `json:"target_x" validate:"required"`
	TargetY      float32 `json:"target_y" validate:"required"`
}

type WorkflowNode struct {
	ID          string   `json:"id" validate:"required"`
	Type        NodeType `json:"type" validate:"required"`
	Initialized bool     `json:"initialized" validate:"required"`
	Position    struct {
		X float32 `json:"x" validate:"required"`
		Y float32 `json:"y" validate:"required"`
	} `json:"position" validate:"required"`
	Definition NodeDefinition `json:"definition" validate:"required"`
}

type NodeDefinition struct {
	Type   TaskType             `json:"type" validate:"required"`
	Fields map[string]NodeField `json:"fields" validate:"required"`
	// TimeoutSec uint16               `json:"timeout_sec" validate:"required,gte=1"`
}

const (
	NodeDefinitionFieldRaybotID  = "raybot_id"
	NodeDefinitionFieldDirection = "direction"
	NodeDefinitionFieldLocation  = "location"
	NodeDefinitionFieldDistance  = "distance"
)

// TODO: NodeField Value can handle generic type
type NodeField struct {
	UseEnv bool    `json:"use_env" validate:"required"`
	Key    *string `json:"key"`
	Value  *string `json:"value"`
}

type ViewPort struct {
	X    float32 `json:"x" validate:"required"`
	Y    float32 `json:"y" validate:"required"`
	Zoom float32 `json:"zoom" validate:"required"`
}

type WorkflowDefinition struct {
	Nodes    []WorkflowNode `json:"nodes" validate:"required"`
	Edges    []WorkflowEdge `json:"edges" validate:"required"`
	Position []float32      `json:"position" validate:"required"`
	ViewPort ViewPort       `json:"view_port" validate:"required"`
	Zoom     float32        `json:"zoom" validate:"required"`
	// Env      map[string]string `json:"env" validate:"required"`
	// TimeoutSec uint64            `json:"timeout_sec" validate:"required,gte=1"`
}
