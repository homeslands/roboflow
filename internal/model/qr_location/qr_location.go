package qrlocation

import (
	"time"

	"github.com/google/uuid"
)

type QRLocation struct {
	ID        string
	Name      string
	QRCode    string
	Metadata  map[string]any
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewQRLocation(name, qrCode string, metadata map[string]any) QRLocation {
	now := time.Now()
	return QRLocation{
		ID:        uuid.NewString(),
		Name:      name,
		QRCode:    qrCode,
		Metadata:  metadata,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
