package raybotcommand

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Type is the type of the raybot command.
type Type string

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (t *Type) UnmarshalText(text []byte) error {
	cmd := Type(text)
	if _, ok := TypeMap[cmd]; !ok {
		return fmt.Errorf("invalid RaybotCommandType: %s", text)
	}
	*t = cmd
	return nil
}

const (
	TypeStop           Type = "STOP"
	TypeMoveForward    Type = "MOVE_FORWARD"
	TypeMoveBackward   Type = "MOVE_BACKWARD"
	TypeMoveToLocation Type = "MOVE_TO_LOCATION"
	TypeOpenBox        Type = "OPEN_BOX"
	TypeCloseBox       Type = "CLOSE_BOX"
	TypeLiftBox        Type = "LIFT_BOX"
	TypeDropBox        Type = "DROP_BOX"
	TypeCheckQRCode    Type = "CHECK_QR"
	TypeWaitGetItem    Type = "WAIT_GET_ITEM"
	TypeScanLocation   Type = "SCAN_LOCATION"
	TypeSpeak          Type = "SPEAK"
)

var TypeMap = map[Type]struct{}{
	TypeStop:           {},
	TypeMoveForward:    {},
	TypeMoveBackward:   {},
	TypeMoveToLocation: {},
	TypeOpenBox:        {},
	TypeCloseBox:       {},
	TypeLiftBox:        {},
	TypeDropBox:        {},
	TypeCheckQRCode:    {},
	TypeWaitGetItem:    {},
	TypeScanLocation:   {},
	TypeSpeak:          {},
}

type Status string

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (s *Status) UnmarshalText(text []byte) error {
	status := Status(text)
	if _, ok := RaybotCommandStatusMap[status]; !ok {
		return fmt.Errorf("invalid RaybotCommandStatus: %s", text)
	}
	*s = status
	return nil
}

const (
	RaybotCommandStatusPending    Status = "PENDING"
	RaybotCommandStatusInProgress Status = "IN_PROGRESS"
	RaybotCommandStatusSucceeded  Status = "SUCCEEDED"
	RaybotCommandStatusFailed     Status = "FAILED"
)

var RaybotCommandStatusMap = map[Status]struct{}{
	RaybotCommandStatusPending:    {},
	RaybotCommandStatusInProgress: {},
	RaybotCommandStatusSucceeded:  {},
	RaybotCommandStatusFailed:     {},
}

type RaybotCommand struct {
	ID          string
	RaybotID    string
	Type        Type
	Status      Status
	Inputs      Inputs
	Outputs     Outputs
	Error       *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CompletedAt *time.Time
}

func NewRaybotCommand(raybotID string, commandType Type, input Inputs) RaybotCommand {
	now := time.Now()
	return RaybotCommand{
		ID:          uuid.NewString(),
		RaybotID:    raybotID,
		Type:        commandType,
		Status:      RaybotCommandStatusPending,
		Inputs:      input,
		Outputs:     NewOutputs([]byte(`{}`)),
		Error:       nil,
		CompletedAt: nil,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
