package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/tuanvumaihuynh/roboflow/internal/raybot/model"
)

type CreateRaybotRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
}

type RaybotResponse struct {
	ID        uuid.UUID          `json:"id" `
	Name      string             `json:"name"`
	Token     string             `json:"token"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	Status    model.RaybotStatus `json:"status"`
}

func MapRaybotToResponse(r *model.Raybot) *RaybotResponse {
	return &RaybotResponse{
		ID:        r.ID,
		Name:      r.Name,
		Token:     r.Token,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
		Status:    r.Status,
	}
}
