// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package api

import (
	"encoding/json"
	"time"

	openapi_types "github.com/oapi-codegen/runtime/types"
)

// CreateQRLocationRequest defines model for CreateQRLocationRequest.
type CreateQRLocationRequest struct {
	Metadata map[string]interface{} `json:"metadata"`
	Name     string                 `json:"name"`
	QrCode   string                 `json:"qrCode"`
}

// CreateRaybotCommandRequest defines model for CreateRaybotCommandRequest.
type CreateRaybotCommandRequest struct {
	Input json.RawMessage   `json:"input"`
	Type  RaybotCommandType `json:"type"`
}

// CreateRaybotRequest defines model for CreateRaybotRequest.
type CreateRaybotRequest struct {
	Name string `json:"name"`
}

// CreateWorkflowRequest defines model for CreateWorkflowRequest.
type CreateWorkflowRequest struct {
	Definition  WorkflowDefinition `json:"definition"`
	Description *string            `json:"description"`
	Name        string             `json:"name"`
}

// Error defines model for Error.
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// NodeDefinition defines model for NodeDefinition.
type NodeDefinition struct {
	Fields map[string]NodeField `json:"fields"`
	Type   TaskType             `json:"type"`
}

// NodeField defines model for NodeField.
type NodeField struct {
	UseEnv bool    `json:"useEnv"`
	Value  *string `json:"value"`
}

// NodeType defines model for NodeType.
type NodeType = string

// PagingQRLocationResponse defines model for PagingQRLocationResponse.
type PagingQRLocationResponse struct {
	Items      []QRLocationResponse `json:"items"`
	TotalItems int64                `json:"totalItems"`
}

// PagingRaybotCommandResponse defines model for PagingRaybotCommandResponse.
type PagingRaybotCommandResponse struct {
	Items      []RaybotCommandResponse `json:"items"`
	TotalItems int64                   `json:"totalItems"`
}

// PagingRaybotResponse defines model for PagingRaybotResponse.
type PagingRaybotResponse struct {
	Items      []RaybotResponse `json:"items"`
	TotalItems int64            `json:"totalItems"`
}

// PagingWorkflowExecutionResponse defines model for PagingWorkflowExecutionResponse.
type PagingWorkflowExecutionResponse struct {
	Items      []WorkflowExecutionResponse `json:"items"`
	TotalItems int64                       `json:"totalItems"`
}

// PagingWorkflowItemResponse defines model for PagingWorkflowItemResponse.
type PagingWorkflowItemResponse struct {
	CreatedAt   time.Time          `json:"createdAt"`
	Description *string            `json:"description"`
	Id          openapi_types.UUID `json:"id"`
	Name        string             `json:"name"`
	UpdatedAt   time.Time          `json:"updatedAt"`
}

// PagingWorkflowResponse defines model for PagingWorkflowResponse.
type PagingWorkflowResponse struct {
	Items      []PagingWorkflowItemResponse `json:"items"`
	TotalItems int64                        `json:"totalItems"`
}

// Position defines model for Position.
type Position struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

// QRLocationResponse defines model for QRLocationResponse.
type QRLocationResponse struct {
	CreatedAt time.Time              `json:"createdAt"`
	Id        openapi_types.UUID     `json:"id"`
	Metadata  map[string]interface{} `json:"metadata"`
	Name      string                 `json:"name"`
	QrCode    string                 `json:"qrCode"`
	UpdatedAt time.Time              `json:"updatedAt"`
}

// RaybotCommandResponse defines model for RaybotCommandResponse.
type RaybotCommandResponse struct {
	CompletedAt *time.Time             `json:"completedAt"`
	CreatedAt   time.Time              `json:"createdAt"`
	Id          openapi_types.UUID     `json:"id"`
	Input       map[string]interface{} `json:"input"`
	Output      map[string]interface{} `json:"output"`
	RaybotId    openapi_types.UUID     `json:"raybotId"`
	Status      RaybotCommandStatus    `json:"status"`
	Type        RaybotCommandType      `json:"type"`
}

// RaybotCommandStatus defines model for RaybotCommandStatus.
type RaybotCommandStatus = string

// RaybotCommandType defines model for RaybotCommandType.
type RaybotCommandType = string

// RaybotResponse defines model for RaybotResponse.
type RaybotResponse struct {
	CreatedAt       time.Time          `json:"createdAt"`
	Id              openapi_types.UUID `json:"id"`
	IpAddress       *string            `json:"ipAddress"`
	LastConnectedAt *time.Time         `json:"lastConnectedAt"`
	Name            string             `json:"name"`
	Status          RaybotStatus       `json:"status"`
	Token           string             `json:"token"`
	UpdatedAt       time.Time          `json:"updatedAt"`
}

// RaybotStatus defines model for RaybotStatus.
type RaybotStatus = string

// RunWorkflowRequest defines model for RunWorkflowRequest.
type RunWorkflowRequest struct {
	Env map[string]string `json:"env"`
}

// RunWorkflowResponse defines model for RunWorkflowResponse.
type RunWorkflowResponse struct {
	WorkflowExecutionId openapi_types.UUID `json:"workflowExecutionId"`
}

// StepResponse defines model for StepResponse.
type StepResponse struct {
	CompletedAt         *time.Time                  `json:"completedAt"`
	Env                 map[string]string           `json:"env"`
	Id                  openapi_types.UUID          `json:"id"`
	Node                WorkflowNode                `json:"node"`
	StartedAt           *time.Time                  `json:"startedAt"`
	Status              WorkflowExecutionStepStatus `json:"status"`
	WorkflowExecutionId openapi_types.UUID          `json:"workflowExecutionId"`
}

// TaskType defines model for TaskType.
type TaskType = string

// UpdateQRLocationRequest defines model for UpdateQRLocationRequest.
type UpdateQRLocationRequest struct {
	Metadata map[string]interface{} `json:"metadata"`
	Name     string                 `json:"name"`
	QrCode   string                 `json:"qrCode"`
}

// UpdateWorkflowRequest defines model for UpdateWorkflowRequest.
type UpdateWorkflowRequest struct {
	Definition  WorkflowDefinition `json:"definition"`
	Description *string            `json:"description"`
	Name        string             `json:"name"`
}

// ViewPort defines model for ViewPort.
type ViewPort struct {
	X    float32 `json:"x"`
	Y    float32 `json:"y"`
	Zoom float32 `json:"zoom"`
}

// WorkflowDefinition defines model for WorkflowDefinition.
type WorkflowDefinition struct {
	Edges    []WorkflowEdge    `json:"edges"`
	Env      map[string]string `json:"env"`
	Nodes    []WorkflowNode    `json:"nodes"`
	Position []float32         `json:"position"`
	Viewport ViewPort          `json:"viewport"`
	Zoom     float32           `json:"zoom"`
}

// WorkflowEdge defines model for WorkflowEdge.
type WorkflowEdge struct {
	Animated     bool    `json:"animated"`
	Id           string  `json:"id"`
	Label        string  `json:"label"`
	Source       string  `json:"source"`
	SourceHandle string  `json:"sourceHandle"`
	SourceX      float32 `json:"sourceX"`
	SourceY      float32 `json:"sourceY"`
	Target       string  `json:"target"`
	TargetHandle string  `json:"targetHandle"`
	TargetX      float32 `json:"targetX"`
	TargetY      float32 `json:"targetY"`
	Type         string  `json:"type"`
}

// WorkflowExecutionResponse defines model for WorkflowExecutionResponse.
type WorkflowExecutionResponse struct {
	CompletedAt *time.Time              `json:"completedAt"`
	CreatedAt   time.Time               `json:"createdAt"`
	Definition  WorkflowDefinition      `json:"definition"`
	Env         map[string]string       `json:"env"`
	Id          openapi_types.UUID      `json:"id"`
	StartedAt   *time.Time              `json:"startedAt"`
	Status      WorkflowExecutionStatus `json:"status"`
	WorkflowId  openapi_types.UUID      `json:"workflowId"`
}

// WorkflowExecutionStatus defines model for WorkflowExecutionStatus.
type WorkflowExecutionStatus = string

// WorkflowExecutionStepStatus defines model for WorkflowExecutionStepStatus.
type WorkflowExecutionStepStatus = string

// WorkflowNode defines model for WorkflowNode.
type WorkflowNode struct {
	Definition  NodeDefinition `json:"definition"`
	Id          string         `json:"id"`
	Initialized bool           `json:"initialized"`
	Position    Position       `json:"position"`
	Type        NodeType       `json:"type"`
}

// WorkflowResponse defines model for WorkflowResponse.
type WorkflowResponse struct {
	CreatedAt   time.Time           `json:"createdAt"`
	Definition  *WorkflowDefinition `json:"definition,omitempty"`
	Description *string             `json:"description"`
	Id          openapi_types.UUID  `json:"id"`
	Name        string              `json:"name"`
	UpdatedAt   time.Time           `json:"updatedAt"`
}

// Page defines model for Page.
type Page = int32

// PageSize defines model for PageSize.
type PageSize = int32

// QRLocationId defines model for QRLocationId.
type QRLocationId = openapi_types.UUID

// RaybotCommandId defines model for RaybotCommandId.
type RaybotCommandId = openapi_types.UUID

// RaybotId defines model for RaybotId.
type RaybotId = openapi_types.UUID

// StepId defines model for StepId.
type StepId = openapi_types.UUID

// WorkflowExecutionId defines model for WorkflowExecutionId.
type WorkflowExecutionId = openapi_types.UUID

// WorkflowId defines model for WorkflowId.
type WorkflowId = openapi_types.UUID

// ListQRLocationsParams defines parameters for ListQRLocations.
type ListQRLocationsParams struct {
	// Page The page number to retrieve (starting from 1).
	Page *Page `form:"page,omitempty" json:"page,omitempty"`

	// PageSize The number of items per page.
	PageSize *int32 `form:"pageSize,omitempty" json:"pageSize,omitempty"`

	// Sort Sort the results by one or more columns.
	//   - Use a column name for ascending order (e.g., created_at).
	//   - Prefix with `-` for descending order (e.g., -created_at).
	//   - Separate multiple columns with a comma (e.g., created_at,-updated_at).
	//
	// Allowed columns: `name`, `qr_code`, `created_at`, `updated_at`.
	Sort *string `form:"sort,omitempty" json:"sort,omitempty"`
}

// ListRaybotsParams defines parameters for ListRaybots.
type ListRaybotsParams struct {
	// Page The page number to retrieve (starting from 1).
	Page *Page `form:"page,omitempty" json:"page,omitempty"`

	// PageSize The number of items per page.
	PageSize *PageSize `form:"pageSize,omitempty" json:"pageSize,omitempty"`

	// Sort Sort the results by one or more columns.
	//   - Use a column name for ascending order (e.g., created_at).
	//   - Prefix with `-` for descending order (e.g., -created_at).
	//   - Separate multiple columns with a comma (e.g., created_at,-updated_at).
	//
	// Allowed columns: `name`, `status`, `last_connected_at`, `created_at`, `updated_at`.
	Sort *string `form:"sort,omitempty" json:"sort,omitempty"`

	// Status Filter by raybot status.
	Status *string `form:"status,omitempty" json:"status,omitempty"`
}

// ListRaybotCommandsParams defines parameters for ListRaybotCommands.
type ListRaybotCommandsParams struct {
	// Page The page number to retrieve (starting from 1).
	Page *Page `form:"page,omitempty" json:"page,omitempty"`

	// PageSize The number of items per page.
	PageSize *PageSize `form:"pageSize,omitempty" json:"pageSize,omitempty"`

	// Sort Sort the results by one or more columns.
	//   - Use a column name for ascending order (e.g., created_at).
	//   - Prefix with `-` for descending order (e.g., -created_at).
	//   - Separate multiple columns with a comma (e.g., created_at,-updated_at).
	//
	// Allowed columns: `status`, `started_at`, `completed_at`.
	Sort *string `form:"sort,omitempty" json:"sort,omitempty"`
}

// ListStepsByWorkflowExecutionIdParams defines parameters for ListStepsByWorkflowExecutionId.
type ListStepsByWorkflowExecutionIdParams struct {
	// Sort Sort the results by one or more columns.
	//   - Use a column name for ascending order (e.g., created_at).
	//   - Prefix with `-` for descending order (e.g., -created_at).
	//   - Separate multiple columns with a comma (e.g., created_at,-updated_at).
	//
	// Allowed columns: `status`, `started_at`, `completed_at`.
	Sort *string `form:"sort,omitempty" json:"sort,omitempty"`
}

// ListWorkflowsParams defines parameters for ListWorkflows.
type ListWorkflowsParams struct {
	// Page The page number to retrieve (starting from 1).
	Page *Page `form:"page,omitempty" json:"page,omitempty"`

	// PageSize The number of items per page.
	PageSize *PageSize `form:"pageSize,omitempty" json:"pageSize,omitempty"`

	// Sort Sort the results by one or more columns.
	//   - Use a column name for ascending order (e.g., created_at).
	//   - Prefix with `-` for descending order (e.g., -created_at).
	//   - Separate multiple columns with a comma (e.g., created_at,-updated_at).
	//
	// Allowed columns: `name`, `started_at`, `completed_at`.
	Sort *string `form:"sort,omitempty" json:"sort,omitempty"`
}

// ListWorkflowExecutionsByWorkflowIdParams defines parameters for ListWorkflowExecutionsByWorkflowId.
type ListWorkflowExecutionsByWorkflowIdParams struct {
	// Page The page number to retrieve (starting from 1).
	Page *Page `form:"page,omitempty" json:"page,omitempty"`

	// PageSize The number of items per page.
	PageSize *PageSize `form:"pageSize,omitempty" json:"pageSize,omitempty"`

	// Sort Sort the results by one or more columns.
	//   - Use a column name for ascending order (e.g., created_at).
	//   - Prefix with `-` for descending order (e.g., -created_at).
	//   - Separate multiple columns with a comma (e.g., created_at,-updated_at).
	//
	// Allowed columns: `status`, `created_at`, `started_at`, `completed_at`.
	Sort *string `form:"sort,omitempty" json:"sort,omitempty"`
}

// CreateQRLocationJSONRequestBody defines body for CreateQRLocation for application/json ContentType.
type CreateQRLocationJSONRequestBody = CreateQRLocationRequest

// UpdateQRLocationByIdJSONRequestBody defines body for UpdateQRLocationById for application/json ContentType.
type UpdateQRLocationByIdJSONRequestBody = UpdateQRLocationRequest

// CreateRaybotJSONRequestBody defines body for CreateRaybot for application/json ContentType.
type CreateRaybotJSONRequestBody = CreateRaybotRequest

// CreateRaybotCommandJSONRequestBody defines body for CreateRaybotCommand for application/json ContentType.
type CreateRaybotCommandJSONRequestBody = CreateRaybotCommandRequest

// CreateWorkflowJSONRequestBody defines body for CreateWorkflow for application/json ContentType.
type CreateWorkflowJSONRequestBody = CreateWorkflowRequest

// UpdateWorkflowByIdJSONRequestBody defines body for UpdateWorkflowById for application/json ContentType.
type UpdateWorkflowByIdJSONRequestBody = UpdateWorkflowRequest

// RunWorkflowByIdJSONRequestBody defines body for RunWorkflowById for application/json ContentType.
type RunWorkflowByIdJSONRequestBody = RunWorkflowRequest
