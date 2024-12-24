package raybot

import (
	"github.com/google/uuid"
	"github.com/tuanvumaihuynh/roboflow/internal/model"
	"github.com/tuanvumaihuynh/roboflow/pkg/validator"
)

type CreateRaybotCommand struct {
	Name string `validate:"required,alphanumspace,min=1,max=100"`
}

func (c CreateRaybotCommand) Validate() error {
	return validator.Validate(c)
}

type DeleteRaybotCommand struct {
	ID uuid.UUID `validate:"required,uuid"`
}

func (c DeleteRaybotCommand) Validate() error {
	return validator.Validate(c)
}

type UpdateStateCommand struct {
	ID    uuid.UUID          `validate:"required,uuid"`
	State model.RaybotStatus `validate:"required"`
}

func (c UpdateStateCommand) Validate() error {
	if err := validator.Validate(c); err != nil {
		return err
	}

	return c.State.Validate()
}
