package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/tuanvumaihuynh/roboflow/internal/command/model"
)

type CreateCommandRequest struct {
	RaybotID uuid.UUID         `json:"raybot_id" binding:"required,uuid"`
	Type     model.CommandType `json:"type" binding:"required"`
	// Input    map[string]interface{} `json:"input" binding:"required"`
}

type CommandResponse struct {
	ID       uuid.UUID           `json:"id"`
	RaybotID uuid.UUID           `json:"raybot_id"`
	Type     model.CommandType   `json:"type"`
	Status   model.CommandStatus `json:"status"`
	// Input     map[string]interface{} `json:"input"`
	// Output    map[string]interface{} `json:"output"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at"`
	// Error     *string   `json:"error"`
}

func MapCommandToResponse(cmd *model.Command) CommandResponse {
	return CommandResponse{
		ID:       cmd.ID,
		RaybotID: cmd.RaybotID,
		Type:     cmd.Type,
		Status:   cmd.Status,
		// Input:    cmd.Input,
		// Output:   cmd.Output,
		CreatedAt:   cmd.CreatedAt,
		CompletedAt: cmd.CompletedAt,
		// Error:     cmd.Error,
	}
}
