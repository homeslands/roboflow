package repository

import (
	"context"

	"github.com/tuanvumaihuynh/roboflow/internal/db/sqldb"
	qrlocation "github.com/tuanvumaihuynh/roboflow/internal/model/qr_location"
	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/sort"
)

type UpdateQRLocationParams struct {
	ID          string
	Name        string
	SetName     bool
	QRCode      string
	SetQRCode   bool
	Metadata    map[string]any
	SetMetadata bool
}

type QRLocationRepository interface {
	// GetQRLocation gets a QRLocation by its ID.
	GetQRLocation(ctx context.Context, db sqldb.SQLDB, id string) (qrlocation.QRLocation, error)

	// ListQRLocations lists all QRLocations.
	ListQRLocations(ctx context.Context, db sqldb.SQLDB, pagingParams paging.Params, sorts []sort.Sort) (paging.List[qrlocation.QRLocation], error)

	// CreateQRLocation creates a new QRLocation.
	CreateQRLocation(ctx context.Context, db sqldb.SQLDB, qrLocation qrlocation.QRLocation) error

	// UpdateQRLocation updates a QRLocation.
	UpdateQRLocation(ctx context.Context, db sqldb.SQLDB, params UpdateQRLocationParams) (qrlocation.QRLocation, error)

	// DeleteQRLocation deletes a QRLocation.
	DeleteQRLocation(ctx context.Context, db sqldb.SQLDB, id string) error
}
