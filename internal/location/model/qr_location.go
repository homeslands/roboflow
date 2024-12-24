package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type QrLocation struct {
	ID        uuid.UUID
	Name      string
	QRCode    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewQrLocation(name, qrCode string) (*QrLocation, error) {
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}
	if qrCode == "" {
		return nil, errors.New("qr code cannot be empty")
	}
	return &QrLocation{
		ID:        uuid.New(),
		Name:      name,
		QRCode:    qrCode,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}
