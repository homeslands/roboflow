package model

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/xsort"
)

type QRLocation struct {
	ID        uuid.UUID
	Name      string
	QRCode    string
	Metadata  map[string]interface{}
	CreatedAt time.Time
	UpdatedAt time.Time
}

type QRLocationRepository interface {
	Get(ctx context.Context, id uuid.UUID) (QRLocation, error)
	List(ctx context.Context, p paging.Params, sorts []xsort.Sort) (*paging.List[QRLocation], error)
	ExistByQRCode(ctx context.Context, qrCode string) (bool, error)
	Create(ctx context.Context, loc QRLocation) error
	Update(ctx context.Context, loc QRLocation) (QRLocation, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
