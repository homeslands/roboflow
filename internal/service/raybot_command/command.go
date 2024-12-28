package raybotcommand

import (
	"github.com/google/uuid"

	"github.com/tuanvumaihuynh/roboflow/internal/model"
	"github.com/tuanvumaihuynh/roboflow/pkg/validator"
)

type CreateRaybotCommandCommand struct {
	RaybotID uuid.UUID               `validate:"required,uuid"`
	Type     model.RaybotCommandType `validate:"required"`
	Input    []byte                  `validate:"omitempty"`
}

func (c *CreateRaybotCommandCommand) Validate() error {
	err := validator.Validate(c)
	if err != nil {
		return err
	}

	if err := c.Type.Validate(); err != nil {
		return err
	}

	return nil
}

type DeleteRaybotCommandCommand struct {
	ID uuid.UUID `validate:"required,uuid"`
}

func (c DeleteRaybotCommandCommand) Validate() error {
	return validator.Validate(c)
}

type SetStatusInProgessCommand struct {
	ID uuid.UUID `validate:"required,uuid"`
}

func (c SetStatusInProgessCommand) Validate() error {
	return validator.Validate(c)
}

type SetStatusSuccessCommand struct {
	ID     uuid.UUID `validate:"required,uuid"`
	Output []byte    `validate:"omitempty"`
}

func (c SetStatusSuccessCommand) Validate() error {
	return validator.Validate(c)
}

type SetStatusFailedCommand struct {
	ID     uuid.UUID `validate:"required,uuid"`
	Output []byte    `validate:"omitempty"`
}

func (c SetStatusFailedCommand) Validate() error {
	return validator.Validate(c)
}
