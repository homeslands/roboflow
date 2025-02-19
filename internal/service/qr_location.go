package service

import (
	"context"

	qrlocation "github.com/tuanvumaihuynh/roboflow/internal/model/qr_location"
	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/sort"
)

type GetQRLocationParams struct {
	ID string `validate:"required,uuid"`
}

type ListQRLocationsParams struct {
	PagingParams paging.Params `validate:"required"`
	Sorts        []sort.Sort   `validate:"sort=name qr_code created_at updated_at"`
}

type CreateQRLocationParams struct {
	Name     string         `validate:"required,alphanumspace,min=1,max=100"`
	QRCode   string         `validate:"required,alphanumspace,min=1,max=100"`
	Metadata map[string]any `validate:"required"`
}

type UpdateQRLocationParams struct {
	ID          string `validate:"required,uuid"`
	Name        string `validate:"required_if=SetName true,omitempty,alphanumspace,min=1,max=100"`
	SetName     bool
	QRCode      string `validate:"required_if=SetQRCode true,omitempty,alphanumspace,min=1,max=100"`
	SetQRCode   bool
	Metadata    map[string]any `validate:"required_if=SetMetadata true,omitempty"`
	SetMetadata bool
}

type DeleteQRLocationParams struct {
	ID string `validate:"required,uuid"`
}

type QRLocationService interface {
	// GetQRLocation gets a QRLocation by its ID.
	GetQRLocation(ctx context.Context, params GetQRLocationParams) (qrlocation.QRLocation, error)

	// ListQRLocations lists all QRLocations.
	ListQRLocations(ctx context.Context, params ListQRLocationsParams) (paging.List[qrlocation.QRLocation], error)

	// CreateQRLocation creates a new QRLocation.
	CreateQRLocation(ctx context.Context, params CreateQRLocationParams) (qrlocation.QRLocation, error)

	// UpdateQRLocation updates a QRLocation.
	UpdateQRLocation(ctx context.Context, params UpdateQRLocationParams) (qrlocation.QRLocation, error)

	// DeleteQRLocation deletes a QRLocation.
	DeleteQRLocation(ctx context.Context, params DeleteQRLocationParams) error
}
