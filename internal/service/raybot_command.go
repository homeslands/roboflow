package service

import (
	"context"
	"time"

	raybotcommand "github.com/tuanvumaihuynh/roboflow/internal/model/raybot_command"
	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/sort"
)

type GetRaybotCommandParams struct {
	ID string `validate:"required,uuid"`
}

type ListRaybotCommandsByRaybotIDParams struct {
	RaybotID     string        `validate:"required,uuid"`
	PagingParams paging.Params `validate:"required"`
	Sorts        []sort.Sort   `validate:"sort=type status completed_at created_at updated_at"`
}

type CreateRaybotCommandParams struct {
	RaybotID string               `validate:"required,uuid"`
	Type     raybotcommand.Type   `validate:"required,enum"`
	Inputs   raybotcommand.Inputs `validate:"required"`
}

type UpdateRaybotCommandParams struct {
	ID             string               `validate:"required,uuid"`
	Status         raybotcommand.Status `validate:"required_if=SetStatus true,omitempty,enum"`
	SetStatus      bool
	Outputs        raybotcommand.Outputs `validate:"required_if=SetOutputs true"`
	SetOutputs     bool
	Error          *string `validate:"required_if=SetError true,omitempty,min=1,max=100"`
	SetError       bool
	CompletedAt    *time.Time `validate:"required_if=SetCompletedAt true"`
	SetCompletedAt bool
}

type DeleteRaybotCommandParams struct {
	ID string `validate:"required,uuid"`
}

type RaybotCommandService interface {
	// GetRaybotCommand gets a raybot command by its ID.
	GetRaybotCommand(ctx context.Context, params GetRaybotCommandParams) (raybotcommand.RaybotCommand, error)

	// ListRaybotCommandsByRaybotID lists all raybot commands by raybot ID.
	ListRaybotCommandsByRaybotID(ctx context.Context, params ListRaybotCommandsByRaybotIDParams) (paging.List[raybotcommand.RaybotCommand], error)

	// CreateRaybotCommand creates a new raybot command.
	CreateRaybotCommand(ctx context.Context, params CreateRaybotCommandParams) (raybotcommand.RaybotCommand, error)

	// UpdateRaybotCommand updates a raybot command.
	UpdateRaybotCommand(ctx context.Context, params UpdateRaybotCommandParams) (raybotcommand.RaybotCommand, error)

	// DeleteRaybotCommand deletes a raybot command.
	DeleteRaybotCommand(ctx context.Context, params DeleteRaybotCommandParams) error
}
