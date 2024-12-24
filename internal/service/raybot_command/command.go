package raybotcommand

import (
	"encoding/json"

	"github.com/google/uuid"

	"github.com/tuanvumaihuynh/roboflow/internal/model"
	"github.com/tuanvumaihuynh/roboflow/pkg/validator"
	"github.com/tuanvumaihuynh/roboflow/pkg/xerrors"
)

type CreateRaybotCommandCommand struct {
	RaybotID uuid.UUID               `validate:"required,uuid"`
	Type     model.RaybotCommandType `validate:"required"`
	Input    any                     `validate:"required"`
}

func (c *CreateRaybotCommandCommand) Validate() error {
	err := validator.Validate(c)
	if err != nil {
		return err
	}

	if err := c.Type.Validate(); err != nil {
		return err
	}

	input, err := validateCommandInput(c.Type, c.Input)
	if err != nil {
		return err
	}
	c.Input = input

	return nil
}

type MoveToLocationInput struct {
	Location  string `json:"location" validate:"required"`
	Direction string `json:"direction" validate:"required,oneof=FORWARD BACKWARD"`
}

func (i MoveToLocationInput) Validate() error {
	return validator.Validate(i)
}

type CheckQRCodeInput struct {
	QRCode string `json:"qr_code" validate:"required"`
}

func (i CheckQRCodeInput) Validate() error {
	return validator.Validate(i)
}

type LiftDropBoxInput struct {
	Distance *int32 `json:"distance" validate:"omitempty gte=300 lte=2000"`
}

func (i LiftDropBoxInput) Validate() error {
	return validator.Validate(i)
}

type DeleteRaybotCommandCommand struct {
	ID uuid.UUID `validate:"required,uuid"`
}

func (c DeleteRaybotCommandCommand) Validate() error {
	return validator.Validate(c)
}

func validateCommandInput(commandType model.RaybotCommandType, input any) (any, error) {
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return nil, xerrors.ThrowInvalidArgument(err, "failed to marshal input")
	}

	switch commandType {
	case model.RaybotCommandTypeLiftBox,
		model.RaybotCommandTypeDropBox:
		var liftDropBoxInput LiftDropBoxInput
		if err := json.Unmarshal(inputBytes, &liftDropBoxInput); err != nil {
			return nil, xerrors.ThrowInvalidArgument(err, "invalid input for this command type")
		}
		return liftDropBoxInput, liftDropBoxInput.Validate()

	case model.RaybotCommandTypeMoveToLocation:
		var moveToLocationInput MoveToLocationInput
		if err := json.Unmarshal(inputBytes, &moveToLocationInput); err != nil {
			return nil, xerrors.ThrowInvalidArgument(err, "invalid input for this command type")
		}
		return moveToLocationInput, moveToLocationInput.Validate()

	case model.RaybotCommandTypeCheckQrCode:
		var checkQRCodeInput CheckQRCodeInput
		if err := json.Unmarshal(inputBytes, &checkQRCodeInput); err != nil {
			return nil, xerrors.ThrowInvalidArgument(err, "invalid input for this command type")
		}
		return checkQRCodeInput, checkQRCodeInput.Validate()

	case model.RaybotCommandTypeStop,
		model.RaybotCommandTypeMoveForward,
		model.RaybotCommandTypeMoveBackward,
		model.RaybotCommandTypeOpenBox,
		model.RaybotCommandTypeCloseBox,
		model.RaybotCommandTypeWaitGetItem:
		// No input required
		if string(inputBytes) != "{}" {
			return nil, xerrors.ThrowInvalidArgument(nil, "no input expected for this command type")
		}
		return nil, nil

	default:
		return nil, xerrors.ThrowInvalidArgument(nil, "unknown command type")
	}
}
