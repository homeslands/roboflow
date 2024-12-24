package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/tuanvumaihuynh/roboflow/internal/location/model"
)

type CreateQRLocationRequest struct {
	Name   string `json:"name" binding:"required,min=1,max=255"`
	QrCode string `json:"qr_code" binding:"required,min=1,max=255"`
}

type UpdateQRLocationRequest struct {
	Name   string `json:"name" binding:"omitempty,min=1,max=255"`
	QrCode string `json:"qr_code" binding:"omitempty,min=1,max=255"`
}

type QRLocationResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	QRCode    string    `json:"qr_code"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func MapQRLocationToResponse(qrLocation *model.QrLocation) QRLocationResponse {
	return QRLocationResponse{
		ID:        qrLocation.ID,
		Name:      qrLocation.Name,
		QRCode:    qrLocation.QRCode,
		CreatedAt: qrLocation.CreatedAt,
		UpdatedAt: qrLocation.UpdatedAt,
	}
}
