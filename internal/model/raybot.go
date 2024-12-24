package model

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/xerrors"
	"github.com/tuanvumaihuynh/roboflow/pkg/xsort"
)

// RaybotStatus represents the state of a raybot
type RaybotStatus string

func (s RaybotStatus) Validate() error {
	switch s {
	case RaybotStatusOffline:
	case RaybotStatusIdle:
	case RaybotStatusBusy:
	default:
		return xerrors.ThrowInvalidArgument(nil, "invalid raybot state")
	}
	return nil
}

const (
	RaybotStatusOffline RaybotStatus = "OFFLINE"
	RaybotStatusIdle    RaybotStatus = "IDLE"
	RaybotStatusBusy    RaybotStatus = "BUSY"
)

type Raybot struct {
	ID              uuid.UUID
	Name            string
	Token           string
	Status          RaybotStatus
	IpAddress       *string
	LastConnectedAt *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type RaybotRepository interface {
	Get(ctx context.Context, id uuid.UUID) (Raybot, error)
	List(ctx context.Context, p paging.Params, sorts []xsort.Sort, state *RaybotStatus) (*paging.List[Raybot], error)
	Create(ctx context.Context, raybot Raybot) error
	Update(ctx context.Context, raybot Raybot) (Raybot, error)
	Delete(ctx context.Context, id uuid.UUID) error

	GetState(ctx context.Context, id uuid.UUID) (RaybotStatus, error)
	UpdateState(ctx context.Context, id uuid.UUID, status RaybotStatus) error
}
