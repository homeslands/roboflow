package model

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/xerrors"
	"github.com/tuanvumaihuynh/roboflow/pkg/xsort"
)

// RaybotCommandStatus represents the state of a command
type RaybotCommandStatus string

func (s RaybotCommandStatus) Validate() error {
	switch s {
	case RaybotCommandStatusPending,
		RaybotCommandStatusInProgress,
		RaybotCommandStatusSuccess,
		RaybotCommandStatusFailed:
	default:
		return xerrors.ThrowInvalidArgument(nil, fmt.Sprintf("invalid command status: %s", s))
	}
	return nil
}

const (
	RaybotCommandStatusPending    RaybotCommandStatus = "PENDING"
	RaybotCommandStatusInProgress RaybotCommandStatus = "IN_PROGRESS"
	RaybotCommandStatusSuccess    RaybotCommandStatus = "SUCCESS"
	RaybotCommandStatusFailed     RaybotCommandStatus = "FAILED"
)

// RaybotCommandType represents the type of a command
type RaybotCommandType string

func (c RaybotCommandType) Validate() error {
	switch c {
	case RaybotCommandTypeStop,
		RaybotCommandTypeMoveForward,
		RaybotCommandTypeMoveBackward,
		RaybotCommandTypeMoveToLocation,
		RaybotCommandTypeOpenBox,
		RaybotCommandTypeCloseBox,
		RaybotCommandTypeLiftBox,
		RaybotCommandTypeDropBox,
		RaybotCommandTypeCheckQrCode,
		RaybotCommandTypeWaitGetItem,
		RaybotCommandTypeScanLocation:
	default:
		return xerrors.ThrowInvalidArgument(nil, fmt.Sprintf("invalid command type: %s", c))
	}
	return nil
}

const (
	RaybotCommandTypeStop           RaybotCommandType = "STOP"
	RaybotCommandTypeMoveForward    RaybotCommandType = "MOVE_FORWARD"
	RaybotCommandTypeMoveBackward   RaybotCommandType = "MOVE_BACKWARD"
	RaybotCommandTypeMoveToLocation RaybotCommandType = "MOVE_TO_LOCATION"
	RaybotCommandTypeOpenBox        RaybotCommandType = "OPEN_BOX"
	RaybotCommandTypeCloseBox       RaybotCommandType = "CLOSE_BOX"
	RaybotCommandTypeLiftBox        RaybotCommandType = "LIFT_BOX"
	RaybotCommandTypeDropBox        RaybotCommandType = "DROP_BOX"
	RaybotCommandTypeCheckQrCode    RaybotCommandType = "CHECK_QR"
	RaybotCommandTypeWaitGetItem    RaybotCommandType = "WAIT_GET_ITEM"

	// Command that has response data
	RaybotCommandTypeScanLocation RaybotCommandType = "SCAN_LOCATION"
)

type MoveToLocationInput struct {
	Location  string `json:"location"`
	Direction string `json:"direction"`
}

type CheckQRCodeInput struct {
	QRCode string `json:"qr_code"`
}

type LiftDropBoxInput struct {
	Distance *int32 `json:"distance"`
}

type RaybotCommand struct {
	RaybotID    uuid.UUID
	ID          uuid.UUID
	Type        RaybotCommandType
	Status      RaybotCommandStatus
	Input       any
	Output      any
	CreatedAt   time.Time
	CompletedAt *time.Time
}

type RaybotCommandRepository interface {
	Get(ctx context.Context, id uuid.UUID) (RaybotCommand, error)
	List(ctx context.Context, raybotId uuid.UUID, p paging.Params, sorts []xsort.Sort) (*paging.List[RaybotCommand], error)
	Create(ctx context.Context, cmd RaybotCommand) error

	// Update updates both the raybot command and its associated raybot's status
	Update(
		ctx context.Context,
		cmdID uuid.UUID,
		raybotStatus RaybotStatus,
		fn func(raybotCmd *RaybotCommand) error,
	) error
	Delete(ctx context.Context, id uuid.UUID) error
}
