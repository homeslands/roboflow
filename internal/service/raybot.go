package service

import (
	"context"
	"time"

	"github.com/tuanvumaihuynh/roboflow/internal/model/raybot"
	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/sort"
)

type GetRaybotParams struct {
	ID string `validate:"required,uuid"`
}

type ListRaybotsParams struct {
	PagingParams paging.Params       `validate:"required"`
	Sorts        []sort.Sort         `validate:"sort=name is_online control_mode ip_address last_connected_at created_at updated_at"`
	IsOnline     *bool               `validate:"omitempty"`
	ControlMode  *raybot.ControlMode `validate:"omitempty,enum"`
}

type CreateRaybotParams struct {
	Name string `validate:"required,alphanumspace,min=1,max=100"`
}

type UpdateRaybotParams struct {
	ID                 string `validate:"required,uuid"`
	Name               string `validate:"required_if=SetName true,omitempty,alphanumspace,min=1,max=100"`
	SetName            bool
	ControlMode        raybot.ControlMode `validate:"required_if=SetControlMode true,omitempty,enum"`
	SetControlMode     bool
	IsOnline           bool `validate:"required_if=SetIsOnline true,omitempty"`
	SetIsOnline        bool
	IPAddress          *string `validate:"required_if=SetIPAddress true,omitempty"`
	SetIPAddress       bool
	LastConnectedAt    *time.Time `validate:"required_if=SetLastConnectedAt true,omitempty"`
	SetLastConnectedAt bool
}

type DeleteRaybotParams struct {
	ID string `validate:"required,uuid"`
}

type RaybotService interface {
	// GetRaybot gets a raybot by its ID.
	GetRaybot(ctx context.Context, params GetRaybotParams) (raybot.Raybot, error)

	// ListRaybots lists all raybots.
	ListRaybots(ctx context.Context, params ListRaybotsParams) (paging.List[raybot.Raybot], error)

	// CreateRaybot creates a new raybot.
	CreateRaybot(ctx context.Context, params CreateRaybotParams) (raybot.Raybot, error)

	// UpdateRaybot updates a raybot.
	UpdateRaybot(ctx context.Context, params UpdateRaybotParams) (raybot.Raybot, error)

	// DeleteRaybot deletes a raybot.
	DeleteRaybot(ctx context.Context, params DeleteRaybotParams) error
}
