package repository

import (
	"context"
	"time"

	"github.com/tuanvumaihuynh/roboflow/internal/db/sqldb"
	raybotcommand "github.com/tuanvumaihuynh/roboflow/internal/model/raybot_command"
	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/sort"
)

type UpdateRaybotCommandParams struct {
	ID             string
	Status         raybotcommand.Status
	SetStatus      bool
	Outputs        raybotcommand.Outputs
	SetOutputs     bool
	Error          *string
	SetError       bool
	CompletedAt    *time.Time
	SetCompletedAt bool
}

type MarkRaybotCommandFailedParams struct {
	RaybotID string
	Error    string
}

type RaybotCommandRepository interface {
	// GetRaybotCommand gets a RaybotCommand by its ID.
	GetRaybotCommand(ctx context.Context, db sqldb.SQLDB, id string) (raybotcommand.RaybotCommand, error)

	// ListRaybotCommandsByRaybotID lists all RaybotCommands by Raybot ID.
	ListRaybotCommandsByRaybotID(ctx context.Context, db sqldb.SQLDB, raybotID string, pagingParams paging.Params, sorts []sort.Sort) (paging.List[raybotcommand.RaybotCommand], error)

	// CreateRaybotCommand creates a new RaybotCommand.
	CreateRaybotCommand(ctx context.Context, db sqldb.SQLDB, raybotCommand raybotcommand.RaybotCommand) error

	// UpdateRaybotCommand updates a RaybotCommand.
	UpdateRaybotCommand(ctx context.Context, db sqldb.SQLDB, params UpdateRaybotCommandParams) (raybotcommand.RaybotCommand, error)

	// DeleteRaybotCommandsByRaybotID deletes a RaybotCommand.
	DeleteRaybotCommandsByRaybotID(ctx context.Context, db sqldb.SQLDB, raybotID string) error

	// MarkRaybotCommandFailed marks multiple RaybotCommands as failed.
	// Filter by Raybot ID.
	MarkRaybotCommandFailed(ctx context.Context, db sqldb.SQLDB, params MarkRaybotCommandFailedParams) error
}
