package repository

import (
	"context"
	"time"

	"github.com/tuanvumaihuynh/roboflow/internal/db/sqldb"
	"github.com/tuanvumaihuynh/roboflow/internal/model/raybot"
	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/sort"
)

type UpdateRaybotParams struct {
	ID                 string
	Name               string
	SetName            bool
	ControlMode        raybot.ControlMode
	SetControlMode     bool
	IsOnline           bool
	SetIsOnline        bool
	IPAddress          *string
	SetIPAddress       bool
	LastConnectedAt    *time.Time
	SetLastConnectedAt bool
}

type RaybotRepository interface {
	// GetRaybot gets a Raybot by its ID.
	GetRaybot(ctx context.Context, db sqldb.SQLDB, id string) (raybot.Raybot, error)

	// ListRaybots lists all Raybots.
	ListRaybots(
		ctx context.Context,
		db sqldb.SQLDB,
		pagingParams paging.Params,
		sorts []sort.Sort,
		isOnline *bool,
		controlMode *raybot.ControlMode,
	) (paging.List[raybot.Raybot], error)

	// CreateRaybot creates a new Raybot.
	CreateRaybot(ctx context.Context, db sqldb.SQLDB, raybot raybot.Raybot) error

	// UpdateRaybot updates a Raybot.
	UpdateRaybot(ctx context.Context, db sqldb.SQLDB, params UpdateRaybotParams) (raybot.Raybot, error)

	// DeleteRaybot deletes a Raybot and all associated raybot commands.
	DeleteRaybot(ctx context.Context, db sqldb.SQLDB, id string) error
}
